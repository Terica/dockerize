#!/bin/bash
echo "building execwdve"
CTR=$(docker run -d -v ${PWD}:/go/src golang:latest go install execwdve)
# wait for build
while test -n "$(docker ps -qf id=${CTR})"; do sleep 1; docker ps -qf id=${CTR}; done
docker logs ${CTR}
docker cp ${CTR}:/go/bin/execwdve ${HOME}/bin
docker rm ${CTR}
