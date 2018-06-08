#!/bin/bash

cd ./bin
./benchmark --threads=10 --requests=1000 --with_interceptor=false
