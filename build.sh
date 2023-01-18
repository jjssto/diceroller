#!/bin/bash
mkdir -p build
DIR=`pwd`
cd src
go mod tidy
go build -o ../build/diceroller ./
cd $DIR
cp -rf ./js ./build/
cp -rf ./css ./build/
cp -rf ./res ./build/
cp -rf ./db ./build/
cp -rf ./templates ./build/
cp -rf ./diceroller.conf ./build/