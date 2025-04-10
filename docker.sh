#!/bin/bash


if ! podman network ls | grep -q rsvp0
then
    podman network create --driver=bridge rsvp0
fi


if ! podman images | grep -q localhost/our-wedding-rsvp
then 
    podman build -t our-wedding-rsvp:latest .
fi


podman run --restart=always \
    --name our-wedding-rsvp --network rsvp0 \
    --tty \
    -p 8080:8080 \
    -v ./data:/workspace/data \
    localhost/our-wedding-rsvp 