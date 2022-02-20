#!/bin/bash

cat /opt/wait-for-it.sh 
/opt/wait-for-it.sh aws:4566 -s -t 10 -- echo "aws is up!"

exec "$@"