#!/bin/bash
set -x
sudo pkill -x openshift
docker ps | awk 'index($NF,"k8s_")==1 { print $1 }' | xargs -l -r docker stop
mount | grep "openshift.local.volumes" | awk '{ print $3}' | xargs -l -r sudo umount
sudo rm -rf openshift.local.*
