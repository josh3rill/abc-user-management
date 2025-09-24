#!/bin/bash
set -e

host="$1"
port="$2"
shift 2
cmd="$@"

echo "\033[1;33m[INFO] On first build, MySQL initialization may take 4-5 minutes. Please be patient!\033[0m"

echo "Waiting for $host:$port to be ready..."

until nc -z "$host" "$port"; do
  echo "$host:$port is unavailable - sleeping for 2 seconds"
  sleep 2
done

echo "$host:$port is up - executing command"
exec $cmd