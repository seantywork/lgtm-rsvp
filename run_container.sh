#!/bin/bash

mkdir -p data

mkdir -p data/media

mkdir -p log

if ! podman network ls | grep -q rsvp0
then
    podman network create --driver=bridge rsvp0
fi


if ! podman images | grep -q localhost/lgtm-rsvp
then 
    podman build -t lgtm-rsvp:latest .
fi


podman run -d --replace --restart=always \
    --name lgtm-rsvp --network rsvp0 \
    --tty \
    -p 8080:8080 \
    -v ./data:/workspace/data \
    -v ./log:/workspace/log \
    -v ./public/images/album:/workspace/public/images/album \
    -v ./config.yaml:/workspace/config.yaml \
    localhost/lgtm-rsvp 