Kash golang in-memory cache
==
## In-memory Cache
### Usage
```go
import "github.com/colindith/kash/store"

s := store.GetShardedMapStore()

_ = s.SetWithTimeout("key1", "timbre+", 1 * time.Minute)

_ = s.Set("key2", "Best Bites!")

v, _ := s.Get("key1")

v, _ = s.Get("key2")

fmt.Println(s.DumpAllJSON())
// {"key1":"timbre+","key2":"Best Bites!"}

s.Delete("key1")

```

### Features
* Implemented with sharded map to reduce time waiting for lock
* Features dumping all data into JSON format
* Support cache key with/without timeout

## Cache TCP Server/Client CLI
Connect to cache storage through TCP protocol.
### Usage
```bash
cd kash && make

// Start the TCP server
./bin/server

// Start the client CLI
./bin/client

Ctl+D to exit
```
![screenshot](/cli.png)

### Features
* Support basic cli function
* Up/down arrow key to find the cmd line history
