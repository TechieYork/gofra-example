#!/bin/bash

cd ./bin
./benchmark --addr=localhost:58888 --threads=100 --requests=2000 --with_interceptor=false
