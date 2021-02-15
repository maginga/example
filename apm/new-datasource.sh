#!/bin/sh

./apm tenant datasource alarm default_nest_01
sleep 5s
 
./apm tenant datasource score default_nest_01
sleep 5s

./apm tenant datasource trace default_nest_01
