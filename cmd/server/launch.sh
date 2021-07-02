#!/bin/bash
set -e

trap 'killall server' SIGINT

cd $(dirname $0)

killall server || true
sleep 0.1

go install -v

server -db-location=beijing.db -shard=Beijing -http-addr=localhost:8011 -config-file=sharding.toml &
server -db-location=beijing-r.db -shard=Beijing -http-addr=localhost:8012 -config-file=sharding.toml -replica &

server -db-location=shanghai.db -shard=Shanghai -http-addr=localhost:8021 -config-file=sharding.toml &
server -db-location=shanghai-r.db -shard=Shanghai -http-addr=localhost:8022 -config-file=sharding.toml -replica &

server -db-location=xian.db -shard=Xian -http-addr=localhost:8031 -config-file=sharding.toml &
server -db-location=xian-r.db -shard=Xian -http-addr=localhost:8032 -config-file=sharding.toml -replica &

server -db-location=hangzhou.db -shard=Hangzhou -http-addr=localhost:8041 -config-file=sharding.toml &
server -db-location=hangzhou-r.db -shard=Hangzhou -http-addr=localhost:8042 -config-file=sharding.toml -replica &

wait
