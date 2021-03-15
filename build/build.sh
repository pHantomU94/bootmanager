#!/bin/bash

echo "#############################################################"
echo "#                 bootmanager install                       #"
echo "# author: hsj                                               #"
echo "#############################################################"
echo ""


red='\033[0;31m'
green="\033[0;32m"

if [ ! x"${USER}" = x"root" ];then
    echo -e ${red}"Please rerun `basename $0` as root. sudo ./build.sh" 
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

echo -e ${green}Success!
