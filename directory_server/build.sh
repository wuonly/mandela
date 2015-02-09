#!/bin/bash


rm -rf dist
mkdir dist

go build main.go
mv main ./dist/directory_server
cp -f config.json ./dist/
cp -f node_entry.json ./dist/



