# Go LRU
A simple [Least-Recently-Used](https://en.wikipedia.org/wiki/Cache_algorithms#Examples) cache for Go. Nothing fancy.

## Usage
Assuming Go LRU is imported as `lru`.

```go
// create a cache that can hold 10 elements and add one
c := lru.New(10)
c.Set("first", "This is the value.")

// fetch it back
v, found := c.Get("first")
fmt.Printf("%v: %v\n", found, v) // > true: This is the value

// add 10 elements; this causes "first" to be evicted as it was the oldest entry
for i := 0; i < 10; i++ {
  c.Set(fmt.Sprintf("%d", i), i)
}

// "first" is not found
v, found = c.Get("first")
fmt.Println(found) // > false
```
