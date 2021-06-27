[![Build Status](https://travis-ci.com/dongnguyenvt/hash-limiter.svg?branch=main)](https://travis-ci.com/dongnguyenvt/hash-limiter) [![codecov](https://codecov.io/gh/dongnguyenvt/hash-limiter/branch/main/graph/badge.svg)](https://codecov.io/gh/dongnguyenvt/hash-limiter)

# hash-limiter
simple lock-free limiter using hash and atomic operations

## Usage
```Golang
import (
    limiter "github.com/dongnguyenvt/hash-limiter"
)

func main() {
    // create a limiter with 1024 buckets and limit of 10 concurrent requests 
    rl, err := limiter.NewLimiter(1024, 10)
    // check error
    client := "some client id"
    resource := "some resource id"
    give, ok := rl.Take(client, resource)
    if !ok {
        // return 429 here
        os.Exit(1)
    }
    defer give()
    // handle some expensive tasks
}
```
