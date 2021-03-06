#!/bin/sh
GOLANG=`which go`
SQLITE=`which sqlite3`
GODEP=`which dep`

if [[ $GOLANG == "" ]]; then
    echo "Please install \`golang\` and run this script again"
    exit 0
fi

if [[ $SQLITE == "" ]]; then
    echo "Please install \`sqlite3\` and run this script again"
    exit 0
fi

if [[ $GODEP == "" ]]; then
    echo "Please install \`dep\` and run this script again"
    exit 0
fi

#restore database
sqlite3 /tmp/ijah.db < ijahDump.sql

#run http server
cd repository/inventory/server/http/main
dep ensure
go run main.go