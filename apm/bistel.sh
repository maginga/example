#!/bin/sh

rm -rf /var/tmp/flink/bistel-org-nest-silence

./apm tenant stream refiner bistel-org-nest-silence metatron-apm-stream-regulation-refinement-1.11.2-1.0-SNAPSHOT-2020_10_26_13_26.jar
sleep 3s

./apm tenant stream fdc bistel-org-nest-silence metatron-apm-stream-model-template-univariate-oos-rules-1.11.2-1.0-SNAPSHOT-2020_11_04_09_45.jar
sleep 3s

./apm tenant stream bae bistel-org-nest-silence metatron-apm-stream-model-template-multivariate-unsupervised-1.11.2-1.0-SNAPSHOT-2020_11_04_09_45.jar
sleep 3s

./apm tenant stream alarm bistel-org-nest-silence metatron-apm-stream-alarm-asset-1.11.1-1.0-SNAPSHOT-2020_10_14_11_27.jar
sleep 3s

./apm tenant stream paramalarm bistel-org-nest-silence metatron-apm-stream-alarm-parameter-1.11.1-1.0-SNAPSHOT-2020_10_19_17_12.jar
sleep 3s
