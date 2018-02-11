#!/bin/sh

# Assuming postgres is installed, and setup for your
# user to have acccess,  this script will create/reset
# the your caraway db with only the testing data.
echo "DROP DATABASE caraway" | psql postgres
psql postgres < dbinit.sql
psql caraway  < data.sql
psql caraway  < bookingdata.sql
