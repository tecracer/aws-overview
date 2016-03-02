#!/bin/sh

source setenv.sh

while true
do
  $AO_EXECUTABLE_PATH/aws-overview -log-file=$AO_LOG_PATH/aws.log
  sleep 180
done
