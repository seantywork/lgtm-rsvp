#!/bin/bash


mkdir -p data

if [  ! -f "data/rsvp.db" ]
then

    sqlite3 data/rsvp.db "VACUUM;"

fi

./rsvp.out


