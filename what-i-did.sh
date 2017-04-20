#!/bin/bash
set -e -x
export KUBECONFIG=openshift.local.config/master/admin.kubeconfig
#docker daemon --insecure-registry 172.30.0.0/16 ##check out /usr/lib/systemd/system/docker.service
#sudo systemctl stop firewalld
sudo systemctl stop firewalld
sudo env "PATH=$PATH" openshift start > openshift.local.log 2>&1 &
sleep 1
sleep 2
sleep 3
sudo chmod +rw $KUBECONFIG
sudo chgrp ilackarms $KUBECONFIG
sudo chown ilackarms $KUBECONFIG
oc login -u system:admin

#manageiq stuff
oadm registry -n default --config=openshift.local.config/master/admin.kubeconfig
sleep 10
oc adm new-project management-infra --description="Management Infrastructure"
oc create sa -n management-infra management-admin
oc create sa -n management-infra inspector-admin
oc create -f - <<API
---
apiVersion: v1
kind: ClusterRole
metadata:
  name: management-infra-admin
rules:
- resources:
  - pods/proxy
  verbs:
  - '*'
API
oc create -f - <<API
apiVersion: v1
kind: ClusterRole
metadata:
  name: hawkular-metrics-admin
rules:
- apiGroups:
  - ""
  resources:
  - hawkular-metrics
  - hawkular-alerts
  verbs:
  - '*'
API
oc adm policy add-role-to-user -n management-infra admin -z management-admin
oc adm policy add-role-to-user -n management-infra management-infra-admin -z management-admin
oc adm policy add-cluster-role-to-user cluster-reader system:serviceaccount:management-infra:management-admin
oc adm policy add-scc-to-user privileged system:serviceaccount:management-infra:management-admin
oc adm policy add-cluster-role-to-user system:image-puller system:serviceaccount:management-infra:inspector-admin
oc adm policy add-scc-to-user privileged system:serviceaccount:management-infra:inspector-admin
oc adm policy add-cluster-role-to-user self-provisioner system:serviceaccount:management-infra:management-admin
oc adm policy add-cluster-role-to-user hawkular-metrics-admin system:serviceaccount:management-infra:management-admin

tail -f openshift.local.log
return

#monitoring stuf
oc project openshift-infra
oc adm policy add-role-to-user view system:serviceaccount:openshift-infra:hawkular -n openshift-infra
oc create -f - <<API
apiVersion: v1
kind: ServiceAccount
metadata:
  name: metrics-deployer
secrets:
- name: metrics-deployer
API
oadm policy add-role-to-user \
    edit system:serviceaccount:openshift-infra:metrics-deployer
oadm policy add-cluster-role-to-user \
    cluster-reader system:serviceaccount:openshift-infra:heapster
oc process -f metrics-deployer.yaml \
    -p USE_PERSISTENT_STORAGE=false \
    -p HAWKULAR_METRICS_HOSTNAME=hawkular-metrics.example.com \
    | oc create -f -
oc secrets new metrics-deployer nothing=/dev/null
