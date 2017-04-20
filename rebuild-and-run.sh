#!/bin/bash
./kill-it.sh && make clean build; rm openshift.local.log; sudo pkill -x openshift; (sudo env "PATH=$PATH" openshift start > openshift.local.log 2>&1 &) && tail -f openshift.local.log
