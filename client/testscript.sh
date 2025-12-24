#!/bin/bash

HOST="localhost"
PORT=8080
DB="data"
COLLECTION="users.json"

NUM_CLIENTS=5

echo "Starting $NUM_CLIENTS parallel clients..."

for i in $(seq 1 $NUM_CLIENTS); do
  {
    sleep $((RANDOM % 2)) # задержка
    echo "$DB $COLLECTION insert '{\"name\":\"user$i\"}'" | nc $HOST $PORT
    sleep $((RANDOM % 2))
    echo "$DB $COLLECTION find '{}'" | nc $HOST $PORT
    sleep $((RANDOM % 2))
    echo "$DB $COLLECTION delete '{\"name\":\"user$i\"}'" | nc $HOST $PORT
  } &
done

wait

echo "All clients finished."
