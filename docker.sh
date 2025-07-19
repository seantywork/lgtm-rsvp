#!/bin/bash

mkdir -p data

mkdir -p data/media

mkdir -p log

if ! podman network ls | grep -q rsvp0
then
    podman network create --driver=bridge rsvp0
fi


if ! podman images | grep -q localhost/our-wedding-rsvp
then 
    podman build -t our-wedding-rsvp:latest .
fi


podman run -d --restart=always \
    --name our-wedding-rsvp --network rsvp0 \
    --tty \
    -p 8080:8080 \
    -v ./data:/workspace/data \
    -v ./log:/workspace/log \
    -v ./public:/workspace/public \
    localhost/our-wedding-rsvp 