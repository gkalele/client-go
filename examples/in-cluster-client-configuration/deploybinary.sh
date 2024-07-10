#!/bin/bash -x

for i in `sudo /local/binaries/kubernetes/kubectl get node | awk '{ print $1 }'`; do rsync -av /tmp/flannel-probe $i:/tmp/ & done
sudo /local/binaries/kubernetes/kubectl rollout restart ds/flannel-probe

