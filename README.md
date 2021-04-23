## Description

a simple distributed key value store in golang.

### Dependency

### Main idea

Using `hash(key)` to get the index, if that is not euqal to current index, then redirecting to the index-matched server

