#!/bin/bash

set -eu

echo "Checking database connection..."

max_wait=30
i=0
until [ $i -ge $max_wait ]
do
  nc -z bbs-go-mysql 3306 && break

  i=$(( i + 1 ))

  echo "$i: Waiting for database..."
  sleep 1
done

if [ $i -eq $max_wait ]
then
  echo "Database connection refused, terminating..."
  exit 1
fi

echo "Database is up, server starting..."

./bbs-go-server