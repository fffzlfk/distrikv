#!/bin/bash
set -e

trap 'killall distrikv' SIGINT

cd $(dirname $0)

killall distrikv || true
sleep 0.1

go install -v

distrikv -db-location=beijing.db -shard=Beijing -http-addr localhost:8080 -config-file=sharding.toml &
distrikv -db-location=shanghai.db -shard=Shanghai -http-addr localhost:8081 -config-file=sharding.toml &
distrikv -db-location=xian.db -shard=Xian -http-addr localhost:8082 -config-file=sharding.toml &

wait