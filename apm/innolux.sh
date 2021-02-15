#!/bin/sh

rm -rf /var/tmp/flink/innolux-org-nest-voice

./apm tenant stream refiner innolux-org-nest-voice metatron-apm-stream-regulation-refinement-1.11.2-1.0-SNAPSHOT-2020_10_26_13_26.jar
sleep 3s

./apm tenant stream bae innolux-org-nest-voice metatron-apm-stream-model-template-multivariate-unsupervised-1.11.2-1.0-SNAPSHOT-2020_11_04_09_45.jar
sleep 3s

./apm tenant stream alarm innolux-org-nest-voice metatron-apm-stream-alarm-asset-1.11.1-1.0-SNAPSHOT-2020_10_14_11_27.jar
sleep 3s
