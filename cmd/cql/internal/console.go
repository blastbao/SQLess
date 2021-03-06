/*
 * Copyright 2016-2018 Kenneth Shaw.
 * Copyright 2018-2019 The CovenantSQL Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package internal

import (
	"bufio"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"time"

	sqlite3 "github.com/SQLess/go-sqlite3-cipher"
	"github.com/xo/dburl"
	"github.com/xo/usql/drivers"
	"github.com/xo/usql/env"
	"github.com/xo/usql/handler"
	"github.com/xo/usql/rline"
	"github.com/xo/usql/text"

	"github.com/SQLess/SQLess/client"
	"github.com/SQLess/SQLess/conf"
	"github.com/SQLess/SQLess/utils"
)

// CmdConsole is cql console command entity.
var CmdConsole = &Command{
	UsageLine: "cql console [common params] [-command sqlcommand] [-out outputfile] [-no-rc true/false] [-single-transaction] [-variable variables] [-explorer explorer_addr] [dsn]",
	Short:     "run a console for interactive sql operation",
	Long: `
Console runs an interactive SQL console for CQL.
e.g.
    cql console cql://4119ef997dedc585bfbcfae00ab6b87b8486fab323a8e107ea1fd4fc4f7eba5c

There is also a -command param for SQL script, and you can add "< file.sql" at end of command for executing a SQL file.
If those params are set, it will run SQL script and exit without staying console mode.
e.g.
    cql console -command "create table test1(test2 int);" cql://4119ef997dedc585bfbcfae00ab6b87b8486fab323a8e107ea1fd4fc4f7eba5c
`,
	Flag:       flag.NewFlagSet("Console params", flag.ExitOnError),
	CommonFlag: flag.NewFlagSet("Common params", flag.ExitOnError),
	DebugFlag:  flag.NewFlagSet("Debug params", flag.ExitOnError),
}

var (
	variables         varsFlag
	dsn               string
	outFile           string
	noRC              bool
	singleTransaction bool
	command           string
)

func init() {
	CmdConsole.Run = runConsole

	addCommonFlags(CmdConsole)
	addConfigFlag(CmdConsole)
	CmdConsole.Flag.Var(&variables, "variable", "Set variable")
	CmdConsole.Flag.StringVar(&outFile, "out", "", "Record stdout to file")
	CmdConsole.Flag.BoolVar(&noRC, "no-rc", false, "Do not read start up file")
	CmdConsole.Flag.BoolVar(&singleTransaction, "single-transaction", false, "Execute as a single transaction (if non-interactive)")
	CmdConsole.Flag.StringVar(&command, "command", "", "Run only single command (SQL or usql internal command) and exit")
	CmdConsole.Flag.StringVar(&explorerAddr, "explorer", "", "Address serve a database chain explorer, e.g. :8546")
}

// SqTime provides a type that will correctly scan the various timestamps
// values stored by the github.com/mattn/go-sqlite3 driver for time.Time
// values, as well as correctly satisfying the sql/driver/Valuer interface.
type SqTime struct {
	time.Time
}

// Value satisfies the Valuer interface.
func (t SqTime) Value() (driver.Value, error) {
	return t.Time, nil
}

// Scan satisfies the Scanner interface.
func (t *SqTime) Scan(v interface{}) error {
	switch x := v.(type) {
	case time.Time:
		t.Time = x
		return nil
	case []byte:
		return t.parse(string(x))
	case string:
		return t.parse(x)
	}

	return fmt.Errorf("cannot convert type %s to time.Time", reflect.TypeOf(v))
}

// parse attempts to parse string s to t.
func (t *SqTime) parse(s string) error {
	if s == "" {
		return nil
	}

	for _, f := range sqlite3.SQLiteTimestampFormats {
		z, err := time.Parse(f, s)
		if err == nil {
			t.Time = z
			return nil
		}
	}

	return errors.New("could not parse time")
}

type varsFlag struct {
	flag.Value
	vars []string
}

func (v *varsFlag) Get() []string {
	return append([]string{}, v.vars...)
}

func (v *varsFlag) String() string {
	return fmt.Sprintf("%#v", v.vars)
}

func (v *varsFlag) Set(value string) error {
	v.vars = append(v.vars, value)
	return nil
}

// UsqlRegister init xo/usql driver
func usqlRegister() {
	// set command name of usql
	text.CommandName = "cqlprotocol"

	// register SQLite3 database
	drivers.Register("sqlite3", drivers.Driver{
		AllowMultilineComments: true,
		ForceParams: drivers.ForceQueryParameters([]string{
			"loc", "auto",
		}),
		Version: func(db drivers.DB) (string, error) {
			var ver string
			err := db.QueryRow(`SELECT sqlite_version()`).Scan(&ver)
			if err != nil {
				return "", err
			}
			return "SQLite3 " + ver, nil
		},
		Err: func(err error) (string, string) {
			if e, ok := err.(sqlite3.Error); ok {
				return strconv.Itoa(int(e.Code)), e.Error()
			}

			code, msg := "", err.Error()
			if e, ok := err.(sqlite3.ErrNo); ok {
				code = strconv.Itoa(int(e))
			}

			return code, msg
		},
		ConvertBytes: func(buf []byte, tfmt string) (string, error) {
			// attempt to convert buf if it matches a time format, and if it
			// does, then return a formatted time string.
			s := string(buf)
			if s != "" && strings.TrimSpace(s) != "" {
				t := new(SqTime)
				if err := t.Scan(buf); err == nil {
					return t.Format(tfmt), nil
				}
			}
			return s, nil
		},
	})

	// register CQL database
	drivers.Register("cqlprotocol", drivers.Driver{
		AllowMultilineComments: true,
		Version: func(db drivers.DB) (string, error) {
			return Version, nil
		},
		Err: func(err error) (string, string) {
			return "", err.Error()
		},
		ConvertBytes: func(buf []byte, tfmt string) (string, error) {
			// attempt to convert buf if it matches a time format, and if it
			// does, then return a formatted time string.
			s := string(buf)
			if s != "" && strings.TrimSpace(s) != "" {
				t := new(SqTime)
				if err := t.Scan(buf); err == nil {
					return t.Format(tfmt), nil
				}
			}
			return s, nil
		},
		RowsAffected: func(sql.Result) (int64, error) {
			return 0, nil
		},
		Open: func(url *dburl.URL) (handler func(driverName, dataSourceName string) (*sql.DB, error), err error) {
			ConsoleLog.Infof("connecting to %#v", url.DSN)

			// wait for database to become ready
			ctx, cancel := context.WithTimeout(context.Background(), waitTxConfirmationMaxDuration)
			defer cancel()
			if err = client.WaitDBCreation(ctx, url.DSN); err != nil {
				return
			}

			return sql.Open, nil
		},
	})

	// register cqlprotocol:// scheme to dburl
	dburl.Register(dburl.Scheme{
		Driver: "cqlprotocol",
		Generator: func(url *dburl.URL) (string, error) {
			return url.String(), nil
		},
	})
}

func run(u *user.User) (err error) {
	// get working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// handle variables
	for _, v := range variables.Get() {
		if i := strings.Index(v, "="); i != -1 {
			env.Set(v[:i], v[i+1:])
		} else {
			env.Unset(v)
		}
	}

	// create input/output
	interactive := command != ""
	l, err := rline.New(interactive, outFile, env.HistoryFile(u))
	if err != nil {
		return err
	}
	defer l.Close()

	// create handler
	h := handler.New(l, u, wd, true)

	// open dsn
	if err = h.Open(dsn); err != nil {
		return err
	}

	// start transaction
	if singleTransaction {
		if h.IO().Interactive() {
			return text.ErrSingleTransactionCannotBeUsedWithInteractiveMode
		}
		if err = h.Begin(); err != nil {
			return err
		}
	}

	// rc file
	if rc := env.RCFile(u); !noRC && rc != "" {
		if err = h.Include(rc, false); err != nil && err != text.ErrNoSuchFileOrDirectory {
			return err
		}
	}

	if command != "" {
		// one liner command
		h.SetSingleLineMode(true)
		h.Reset([]rune(command))
		if err = h.Run(); err != nil && err != io.EOF {
			ConsoleLog.WithError(err).Error("run command failed")
			SetExitStatus(1)
			return
		}
	} else {
		// interactive
		if err = h.Run(); err != nil {
			return
		}

	}

	// commit
	if singleTransaction {
		return h.Commit()
	}

	return nil
}

// runConsole runs a console for sql operation in command line.
func runConsole(cmd *Command, args []string) {
	configFile = utils.HomeDirExpand(configFile)

	commonFlagsInit(cmd)

	var err error
	// load config
	if conf.GConf, err = conf.LoadConfig(configFile); err != nil {
		ConsoleLog.WithError(err).Error("load config file failed")
		SetExitStatus(1)
		ExitIfErrors()
	}

	if len(args) == 1 {
		dsn = args[0]
	}

	if dsn == "" {
		dsnArray := loadDSN()
		if len(dsnArray) > 0 {
			//Print dsn list
			fmt.Printf("Found local stored dsn list: \n")
			for i := 0; i < len(dsnArray); i++ {
				fmt.Printf("%v: %v\n", i, dsnArray[i])
			}
			fmt.Println("Which would you like to connect? (press Enter for default 0):")

			//Read from terminal
			reader := bufio.NewReader(os.Stdin)
			t, err := reader.ReadString('\n')
			t = strings.Trim(t, "\n")
			if err != nil {
				ConsoleLog.WithError(err).Error("unexpected error")
				SetExitStatus(1)
				Exit()
			}

			var choice int
			if t == "" {
				choice = 0
			} else {
				choice, err = strconv.Atoi(t)
				if err != nil || choice >= len(dsnArray) || choice < 0 {
					ConsoleLog.Error("invalid choice number")
					SetExitStatus(1)
					Exit()
				}
			}

			//Set dsn
			dsn = dsnArray[choice]
		} else {
			ConsoleLog.Error("neither local dsn storage exists nor a dsn string present")
			SetExitStatus(1)
			printCommandHelp(cmd)
			Exit()
		}
	}

	configInit()

	usqlRegister()

	var (
		curUser   *user.User
		available = drivers.Available()
	)
	if st, err := os.Stat("/.dockerenv"); err == nil && !st.IsDir() {
		// in docker, fake user
		var wd string
		if wd, err = os.Getwd(); err != nil {
			ConsoleLog.WithError(err).Error("get working directory failed")
			SetExitStatus(1)
			return
		}
		curUser = &user.User{
			Uid:      "0",
			Gid:      "0",
			Username: "docker",
			Name:     "docker",
			HomeDir:  wd,
		}
	} else {
		if curUser, err = user.Current(); err != nil {
			ConsoleLog.WithError(err).Error("get current user failed")
			SetExitStatus(1)
			return
		}
	}

	if explorerAddr != "" {
		cancelFunc := startExplorerServer(explorerAddr)
		defer cancelFunc()
	}

	// run
	err = run(curUser)
	ExitIfErrors()
	if err != nil && err != io.EOF && err != rline.ErrInterrupt {
		ConsoleLog.WithError(err).Error("run cli error")

		if e, ok := err.(*drivers.Error); ok && e.Err == text.ErrDriverNotAvailable {
			bindings := make([]string, 0, len(available))
			for name := range available {
				bindings = append(bindings, name)
			}
			ConsoleLog.Infof("available drivers are: %#v", bindings)
		}
		SetExitStatus(1)
		return
	}

	if explorerAddr != "" {
		ConsoleLog.Printf("Ctrl + C to stop background server on %s\n", explorerAddr)
		<-utils.WaitForExit()
	}
}
