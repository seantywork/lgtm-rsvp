#!/bin/bash


mkdir -p data

mkdir -p data/media

mkdir -p log


if [  ! -f "data/rsvp.db" ]
then

    sqlite3 data/rsvp.db "VACUUM;"

fi

datestr=$(date '+%Y-%m-%d-%H-%M-%S')

./rsvp.out >> "log/$datestr" 2>&1


