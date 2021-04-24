## Description

a simple distributed key value store in golang.

### Dependency

### Main idea

#### SetKey & GetKey
Using `hash(key)` to get the index, if that is not euqal to current index, then redirecting to the index-matched server


#### Replication

Created a Queue on the Masters, this queue stores key-values that have not yet been written to replicas, Slaves loop get datas from the queue and delete them from Master.

### Configuration

`sharding.toml`

### Use

```bash
./launch.sh
```
