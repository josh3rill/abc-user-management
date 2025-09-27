#!/bin/bash
set -e

host="$1"
port="$2"
shift 2
cmd="$@"

echo "Waiting for $host:$port to be ready..."

until nc -z "$host" "$port"; do
  echo "$host:$port is unavailable - sleeping for 2 seconds"
  sleep 2
done

echo "$host:$port is up - executing command"
exec $cmd