package lru

import (
  "sync"
  "strconv"
  "testing"
  "math/rand"
  "github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
  var v interface{}
  var ok bool
  c := New(3)
  
  c.Set(1, "1")
  v, ok = c.Get(1)
  assert.Equal(t, "1", v)
  c.Set(2, "2")
  v, ok = c.Get(2)
  assert.Equal(t, "2", v)
  c.Set(3, "3")
  v, ok = c.Get(3)
  assert.Equal(t, "3", v)
  
  c.Set(4, "4")
  v, ok = c.Get(4)
  assert.Equal(t, "4", v)
  _, ok = c.Get(1)
  assert.Equal(t, false, ok)
  
  c.Set(5, "5")
  v, ok = c.Get(5)
  assert.Equal(t, "5", v)
  _, ok = c.Get(2)
  assert.Equal(t, false, ok)
  
  c.Delete(3)
  _, ok = c.Get(3)
  assert.Equal(t, false, ok)
  
  c.Set(6, "6")
  v, ok = c.Get(6)
  assert.Equal(t, "6", v)
  v, ok = c.Get(4)
  assert.Equal(t, "4", v)
  
}

func TestConcurrent(t *testing.T) {
  var waiter sync.WaitGroup
  c := New(3)
  
  max := 100000
  sem := make(chan struct{}, 25)
  for i := 0; i < max; i++ {
    sem <- struct{}{}
    waiter.Add(1)
    go func(){
      defer func(){ <-sem; waiter.Done() }()
      c.Set(i, strconv.FormatInt(int64(i), 10))
      for j := 0; j < 10; j++ {
        c.Get(rand.Int() % i) // just do an access
      }
    }()
  }
  
  waiter.Wait()
}
