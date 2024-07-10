#!/bin/bash

ls -lh /binaries
cp /binaries/flannel-probe /tmp/
/tmp/flannel-probe | tee /local/logs/tetration/flannel-probe/current


