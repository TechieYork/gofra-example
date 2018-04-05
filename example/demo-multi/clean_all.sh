#!/bin/bash

cd ./default
./clean.sh

cd ../serviceA
./clean.sh

cd ../serviceB
./clean.sh

cd ../serviceC
./clean.sh

cd ../serviceD
./clean.sh
