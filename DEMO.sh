#!/bin/bash

# usql: https://github.com/xo/usql

CONF=conf/testnet/config.yaml

DSN_FILE=conf/testnet/.dsn
if [ ! -f $DSN_FILE ]
then
bin/cql create -config=$CONF \
    -db-node 1 -wait-tx-confirm	|| exit 1
fi

DSN="`cat $DSN_FILE`"

echo
cat sample-sql-command.txt
echo

set -x

bin/cql console -config=$CONF $DSN
bin/cql explorer -config=$CONF localhost:8081
bin/cql wallet -config=$CONF
