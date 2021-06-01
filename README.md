<h1 align="center">Welcome to distributed-key-value 👋</h1>
<p>
  <a href="https://twitter.com/fffzlfk" target="_blank">
    <img alt="Twitter: fffzlfk" src="https://img.shields.io/twitter/follow/fffzlfk.svg?style=social" />
  </a>
</p>

> A simple distributed key value store in golang

## Main idea

### SetKey & GetKey

Using hash(key) to get the index, if that is not euqal to current index, then redirecting to the index-matched server

### Replication
Created a Queue on the Masters, this queue stores key-values that have not yet been written to replicas, Slaves loop get datas from the queue and delete them from Master.

## Usage

```sh
./launsh.sh
```

### Configuration

[sharding.toml](./sharding.toml)

## Author

👤 **fffzlfk**

* Website: [fffzlfk.netlify.app](https://fffzlfk.netlify.app)
* Twitter: [@fffzlfk](https://twitter.com/fffzlfk)
* Github: [@fffzlfk](https://github.com/fffzlfk)

## Show your support

Give a ⭐️ if this project helped you!

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_
