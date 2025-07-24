#!/bin/bash

print_help(){
    echo "help: print help message"
    echo "flog: use file log"
}

mkdir -p data

mkdir -p data/media


if [  ! -f "data/rsvp.db" ]
then

    sqlite3 data/rsvp.db "VACUUM;"

fi

echo "running"

if [ ! -z "${1}" ]
then 
    if [ "$1" == "flog" ]
    then
        echo "run: use file log"
        currdate=$(date '+%Y-%m-%d-%H-%M-%S')
        ./rsvp.out >> "data/$currdate" 2>&1
    elif [ "$1" == "help" ]
    then 
        print_help
        exit 0
    else
        echo "invalid arg: $1"
        print_help
        exit 1
    fi
else 
    ./rsvp.out
fi



