#!/bin/bash

cd ./default
./stop.sh

cd ../serviceA
./stop.sh

cd ../serviceB
./stop.sh

cd ../serviceC
./stop.sh

cd ../serviceD
./stop.sh
