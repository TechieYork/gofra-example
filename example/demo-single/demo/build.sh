#!/bin/bash

cd ./src
go build -o demo

cd ../test
go build
