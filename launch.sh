#!/bin/bash
set -e

trap 'killall distrikv' SIGINT

cd $(dirname $0)

killall distrikv || true
sleep 0.1

go install -v

distrikv -db-location=beijing.db -shard=Beijing -http-addr=127.0.0.2:8080 -config-file=sharding.toml &
distrikv -db-location=beijing-r.db -shard=Beijing -http-addr=127.0.0.22:8080 -config-file=sharding.toml -replica &

distrikv -db-location=shanghai.db -shard=Shanghai -http-addr=127.0.0.3:8080 -config-file=sharding.toml &
distrikv -db-location=shanghai-r.db -shard=Shanghai -http-addr=127.0.0.33:8080 -config-file=sharding.toml -replica &

distrikv -db-location=xian.db -shard=Xian -http-addr=127.0.0.4:8080 -config-file=sharding.toml &
distrikv -db-location=xian-r.db -shard=Xian -http-addr=127.0.0.44:8080 -config-file=sharding.toml -replica &

distrikv -db-location=hangzhou.db -shard=Hangzhou -http-addr=127.0.0.5:8080 -config-file=sharding.toml &
distrikv -db-location=hangzhou-r.db -shard=Hangzhou -http-addr=127.0.0.55:8080 -config-file=sharding.toml -replica &

wait