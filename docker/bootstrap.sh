#!/bin/sh

set -x

#start basic redis on default port
echo never | sudo tee /sys/kernel/mm/transparent_hugepage/enabled
sudo docker run --rm -p 6379:6379 -d --name redis redis
