# Mango Cache
[![GoDoc](https://godoc.org/github.com/goburrow/cache?status.svg)](https://godoc.org/github.com/goburrow/cache)
[![Build Status](https://travis-ci.org/goburrow/cache.svg?branch=master)](https://travis-ci.org/goburrow/cache)

Partial implementations of [Guava Cache](https://github.com/google/guava) in Go.

Supported cache replacement policies:

- LRU
- Segmented LRU (default)
- TinyLFU (experimental)

The TinyLFU implementation is inspired by
[Caffeine](https://github.com/ben-manes/caffeine) by Ben Manes and
[go-tinylfu](https://github.com/dgryski/go-tinylfu) by Damian Gryski.

## Download

```
go get -u github.com/goburrow/cache
```

## Example

```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/goburrow/cache"
)

func main() {
	load := func(k cache.Key) (cache.Value, error) {
		return fmt.Sprintf("%d", k), nil
	}
	// Create a new cache
	c := cache.NewLoadingCache(load,
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(10*time.Second),
		cache.WithRefreshAfterWrite(60*time.Second),
	)

	getTicker := time.Tick(10 * time.Millisecond)
	reportTicker := time.Tick(1 * time.Second)
	for {
		select {
		case <-getTicker:
			_, _ = c.Get(rand.Intn(2000))
		case <-reportTicker:
			st := cache.Stats{}
			c.Stats(&st)
			fmt.Printf("%+v\n", st)
		}
	}
}
```

## Performance

See [traces](traces/)

![report](traces/report.png)
