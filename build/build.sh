#!/bin/bash

if [ ! x"${USER}" = x"root" ];then
    echo "Please rerun `basename $0` as root. sudo ./build.sh" 
    exit 1
fi


BIN_FILE="bootmanager"
BIN_DIR="/usr/bin"
CONFIGURE_FILE="config.json"
CONFIGURE_DIR="/usr/local/bootmanager"


# sudo cp bootmanager $BIN_DIR

cp $BIN_FILE $BIN_DIR

if [ ! -d $CONFIGURE_DIR ]; then
    mkdir -p $CONFIGURE_DIR
fi

cp $CONFIGURE_FILE $CONFIGURE_DIR
