default: all

version := 0.1-alpha

test_tags := $(tags) testbinary
test_flags := -coverpkg github.com/SQLess/SQLess/... -cover -race -c

ldflags_role_client := -X main.version=$(version) -X github.com/SQLess/SQLess/conf.RoleTag=C
ldflags_role_client_simple_log := $(ldflags_role_client) -X github.com/SQLess/SQLess/utils/log.SimpleLog=Y

GOTEST := CGO_ENABLED=1 go test $(test_flags) -tags "$(test_tags)"
GOBUILD := CGO_ENABLED=1 go build -tags "$(tags)"

bin/cql.test: gen-stamp
	$(GOTEST) \
		-ldflags "$(ldflags_role_client)" \
		-o bin/cql.test \
		github.com/SQLess/SQLess/cmd/cql

bin/cql: gen-stamp
	$(GOBUILD) \
		-ldflags "$(ldflags_role_client_simple_log)" \
		-o bin/cql \
		github.com/SQLess/SQLess/cmd/cql

gen-stamp: bin/hsp
	PATH=`pwd`/bin:$(PATH) go generate ./...
	touch $@

bin/hsp:
	go get -u github.com/SQLess/HashStablePack/hsp
	go build \
		-o bin/hsp \
		github.com/SQLess/HashStablePack/hsp

client: bin/cql bin/cql.test

all: client

clean:
	find . -name '*_gen*.go' -exec rm -vf {} \;
	$(RM) bin/cql* bin/hsp
	$(RM) *.cover.out
	$(RM) coverage.txt
	$(RM) gen-stamp

testnet-clean:
	$(RM) conf/testnet/.dsn
	$(RM) conf/testnet/observer.db3
	$(RM) conf/testnet/public.keystore*
	$(RM) conf/testnet/observer.db3

.PHONY: bin/cql bin/cql.test
