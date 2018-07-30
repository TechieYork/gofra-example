#!/bin/bash

cd ./bin
./benchmark --addr=localhost:58888 --threads=50 --requests=20000 --with_interceptor=false
#./benchmark --addr=localhost:58888 --threads=50 --requests=20000 --with_interceptor=true
