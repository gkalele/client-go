#!/bin/bash -x

for i in `sudo /local/binaries/kubernetes/kubectl get node | awk '{ print $1 }'`;  do
    rsync -av /tmp/probe.tar $i:/tmp/
done
for i in `sudo /local/binaries/kubernetes/kubectl get node | awk '{ print $1 }'`;  do
    ssh $i "sudo ctr -n k8s.io image import /tmp/probe.tar" & done
sudo /local/binaries/kubernetes/kubectl rollout restart ds/flannel-probe
