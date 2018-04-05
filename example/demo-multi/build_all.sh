#!/bin/bash

cd ./default
./build.sh

cd ../serviceA
./build.sh

cd ../serviceB
./build.sh

cd ../serviceC
./build.sh

cd ../serviceD
./build.sh
