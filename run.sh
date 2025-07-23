#!/bin/bash


mkdir -p data

mkdir -p data/media


if [  ! -f "data/rsvp.db" ]
then

    sqlite3 data/rsvp.db "VACUUM;"

fi

if [ "$1" == "filelog" ]
then
    currdate=$(date '+%Y-%m-%d-%H-%M-%S')
    ./rsvp.out >> "data/$currdate" 2>&1
else 
    ./rsvp.out
fi



