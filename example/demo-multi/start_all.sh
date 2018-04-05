#!/bin/bash

cd ./default
./start.sh

cd ../serviceA
./start.sh

cd ../serviceB
./start.sh

cd ../serviceC
./start.sh

cd ../serviceD
./start.sh
