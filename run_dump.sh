#!/bin/bash


sqlite3 data/rsvp.db .dump > rsvpdump.sql


sed -i -e 's/PRAGMA foreign_keys=OFF;//g' rsvpdump.sql
sed -i -e 's/BEGIN TRANSACTION;//g' rsvpdump.sql
sed -i -e 's/INTEGER PRIMARY KEY/INT NOT NULL AUTO_INCREMENT/g' rsvpdump.sql
sed -i -e 's/pw TEXT/pw TEXT,\n    PRIMARY KEY(admin_id)/g' rsvpdump.sql
sed -i -e '0,/content TEXT/ s//content TEXT,\n    PRIMARY KEY(story_id)/' rsvpdump.sql
sed -i -e 's/timestamp_approved TEXT/timestamp_approved TEXT,\n    PRIMARY KEY(comment_id)/g' rsvpdump.sql
