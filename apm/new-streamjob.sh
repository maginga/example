#!/bin/sh

rm -rf /var/tmp/flink/default_nest_01

./apm tenant stream refiner default_nest_01 metatron-apm-stream-regulation-refinement-1.11.2-1.0-SNAPSHOT-2021_02_05_10_56.jar
sleep 3s

./apm tenant stream fdc default_nest_01 metatron-apm-stream-model-template-univariate-oos-rules-1.11.2-1.0-SNAPSHOT-2020_11_04_09_45.jar
sleep 3s

./apm tenant stream bae default_nest_01 metatron-apm-stream-model-template-multivariate-unsupervised-1.11.2-1.0-SNAPSHOT-2020_11_04_09_45.jar
sleep 3s

./apm tenant stream alarm default_nest_01 metatron-apm-stream-alarm-asset-1.11.1-1.0-SNAPSHOT-2020_10_14_11_27.jar
sleep 3s

./apm tenant stream paramalarm default_nest_01 metatron-apm-stream-alarm-parameter-1.11.1-1.0-SNAPSHOT-2020_10_19_17_12.jar
sleep 3s
