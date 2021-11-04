#!/bin/bash
echo "running mongo docker image"

docker run -it --rm --name treehub-mongo \
    -p 8080:8080 \
   mrshuvava/treehub:0.0.6__linux_amd64
