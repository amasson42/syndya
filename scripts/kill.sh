#!/bin/sh

UNAME_S=`uname -s`
LISTEN_PORT=8080

if [ "$UNAME_S" = "Linux" ]; then
    sudo kill -9 `sudo lsof -t -i:$LISTEN_PORT` 2> /dev/null || :
fi
if [ "$UNAME_S" = "Darwin" ]; then
    npx kill-port $LISTEN_PORT 2> /dev/null || :
fi
