#!/usr/bin/env bash

USER="user1"
PASS="pass1"
DB="baseball"
PORT="5433"

sudo -u postgres dropdb -p $PORT --if-exists $DB
sudo -u postgres createdb -p $PORT $DB

sudo -u postgres dropuser -p $PORT --if-exists $USER
sudo -u postgres createuser -p $PORT $USER

sudo -u postgres psql -p $PORT -c "ALTER USER user1 WITH encrypted password '$PASS'"

PGPASSWORD=$PASS psql -U $USER -d $DB -p $PORT -h localhost -w -q -f ./initial.sql
