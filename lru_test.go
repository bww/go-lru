package lru

import (
  "testing"
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
