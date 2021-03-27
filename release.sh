#!/bin/bash
go build .
cp bootmanager ./build/bootmanager/
tar -czvf ./build/bootmanager.tar.gz ./build/bootmanager