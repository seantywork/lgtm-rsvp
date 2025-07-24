#!/bin/bash


print_help(){
    echo "help: print help message"
    echo "flog: use file log"
    echo "logf: follow logs when not using file logging"
    echo "stop: stop container"
    echo "remove: remove container image"
}


mkdir -p data

mkdir -p data/media


if ! podman network ls | grep -q rsvp0
then
    podman network create --driver=bridge rsvp0
fi


if ! podman images | grep -q localhost/lgtm-rsvp
then 
    podman build -t lgtm-rsvp:latest .
fi

FLOG=""

if [ ! -z "${1}" ]
then 

    if [ "$1" == "flog" ]
    then    
        echo "using file logging"
        FLOG="$1"
    elif [ "$1" == "logf" ]
    then 
        podman logs --follow lgtm-rsvp
        exit 0
    elif [ "$1" == "stop" ]
    then 
        podman stop lgtm-rsvp 
        exit 0
    elif [ "$1" == "remove" ]
    then
        podman image rm -f lgtm-rsvp:latest
        exit 0 
    elif [ "$1" == "help" ]
    then 
        print_help
        exit 0
    else 
        echo "invalid command arg: $1"
        print_help
        exit 1
    fi 

else 

    echo "not using filelog"

fi

podman run -d --replace --restart=always \
    --name lgtm-rsvp --network rsvp0 \
    --tty \
    -p 8080:8080 \
    -v ./data:/workspace/data \
    -v ./public/images/album:/workspace/public/images/album \
    -v ./config:/workspace/config \
    localhost/lgtm-rsvp \
    $FLOG
