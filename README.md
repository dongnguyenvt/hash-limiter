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
    ok, give := rl.Take(client, resource)
    if !ok {
        // return 429 here
        os.Exit(1)
    }
    defer give()
    // handle some expensive tasks
}
```
