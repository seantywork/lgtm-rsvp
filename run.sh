#!/bin/bash


mkdir -p data

mkdir -p data/media


if [  ! -f "data/rsvp.db" ]
then

    sqlite3 data/rsvp.db "VACUUM;"

fi

./rsvp.out


