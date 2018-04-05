#!/bin/bash

cd ./src
go build -o default

cd ../test
go build
