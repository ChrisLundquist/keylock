### Keylock: Keys mapped to locks

Create a fixed set of mutexes to which keys map for fine grain locking.

This is probably not very idiomatic go code.

##### Why did you need this?

When using a wonderful [golang lru](https://github.com/hashicorp/golang-lru/) library
I needed fine grain locking for complex merges.
I'm using the above to help batch database writes.


Something like this:
```go
var write_batcher, _ = lru.NewWithEvict(4096, on_evict)
var event_mutex = &sync.Mutex{}

func on_evict(key interface{}, value interface{}) {
    event, _ := value.(Event)
    event.Save()
}

type Event struct {
    Counter int
    // ... Other information to merge
}

func ProcessEvent(event *Event) {
    // ...
    cache_key := event.CacheKey()
    event_mutex.Lock()
    defer event_mutex.Unlock()
    value, found := write_batcher.Get(cache_key)
    if found {
        cached_event := value.(Event)
        cached_event.Counter += event.Counter
        // NOTE: It is possible that the key was evicted since we got our lock
        // This will inflate our counter and duplicate information
        write_batcher.Add(cache_key, cached_event)
    } else {
        // Add it fresh
        write_batcher.Add(cache_key, event)
    }

}
```

The above quickly becomes constrained by `event_mutex`
While we do need to prevent multiple goroutines from working on the same key
Having goroutines working on different keys is acceptable.

In the above example we could replace `event_mutex.Lock()` with `event_mutex.Lock(cache_key)`
and allow for higher concurrency.

We also don't want to allow for an unbounded number of mutexes, as an attacker could generate a large number of events
and create abnormal memory usage.
