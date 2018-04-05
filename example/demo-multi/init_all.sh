#!/bin/bash

cd ./default
./init.sh

cd ../serviceA
./init.sh

cd ../serviceB
./init.sh

cd ../serviceC
./init.sh

cd ../serviceD
./init.sh

cd ..
cp bak/main.go default/test/main.go
cp bak/AddUser.go default/src/handler/UserService/AddUser.go
cp bak/AddAge.go serviceB/src/handler/AgeService/AddAge.go
