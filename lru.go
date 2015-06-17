// 
// Copyright (c) 2015 Brian William Wolter, All rights reserved.
// A simple Least-Recently-Used cache implementation
// 
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
// 
//   * Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
// 
//   * Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//     
//   * Neither the names of Brian William Wolter, Wolter Group New York, nor the
//     names of its contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//     
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
// OF THE POSSIBILITY OF SUCH DAMAGE.
// 

package lru

import (
  "fmt"
  "sync"
  "container/list"
)

/**
 * Cache element
 */
type cacheElement struct {
  key   string
  value interface{}
  elem  *list.Element
}

/**
 * A key-value pair
 */
type KeyValue struct {
  Key   string
  Value interface{}
}

/**
 * A LRU memory cache
 */
type Cache struct {
  sync.RWMutex
  lru       *list.List
  elem      map[string]*cacheElement
  evicted   chan KeyValue
  limit     int
}

/**
 * Create a cache
 */
func NewCache(limit int) *Cache {
  return &Cache{sync.RWMutex{}, list.New(), make(map[string]*cacheElement), nil, limit}
}

/**
 * Obtain the evicted element channel
 */
func (c *Cache) Evicted() <-chan KeyValue {
  c.Lock()
  defer c.Unlock()
  return c.evictedChannel(64)
}

/**
 * Obtain the evicted element channel. The channal is created if it does
 * not yet exist. This method must be externally synchronized.
 */
func (c *Cache) evictedChannel(backlog int) chan KeyValue {
  if c.evicted == nil {
    c.evicted = make(chan KeyValue, backlog)
  }
  return c.evicted
}

/**
 * Display this cache's contents
 */
func (c *Cache) Show() {
  c.Lock()
  defer c.Unlock()
  i := 0
  for k, v := range c.elem {
    fmt.Printf("[%d] %v -> %v\n", i, k, v)
    i++
  }
  i = 0
  for e := c.lru.Front(); e != nil; e = e.Next() {
    fmt.Printf("[%d] %v\n", i, e.Value)
    i++
  }
}

/**
 * Obtain the number of elements
 */
func (c *Cache) Count() int {
  c.RLock()
  defer c.RUnlock()
  return len(c.elem)
}

/**
 * Iterate over elements
 */
func (c *Cache) Iter(f func(string, interface{})) {
  c.RLock()
  defer c.RUnlock()
  for k, v := range c.elem {
    f(k, v.value)
  }
}

/**
 * Get a value
 */
func (c *Cache) Get(key string) (interface{}, bool) {
  c.RLock()
  defer c.RUnlock()
  v, ok := c.elem[key]
  if ok {
    c.lru.MoveToFront(v.elem)
    return v.value, true
  }else{
    return nil, false
  }
}

/**
 * Set a value
 */
func (c *Cache) Set(key string, value interface{}) {
  c.Lock()
  defer c.Unlock()
  v, ok := c.elem[key]
  if ok {
    v.value = value
    c.lru.MoveToFront(v.elem)
  }else{
    if c.limit > 0 {
      for e := c.lru.Back(); len(c.elem) + 1 > c.limit; e = e.Prev() {
        key := e.Value.(string)
        if c.evicted != nil {
          c.evicted <- KeyValue{key, c.elem[key]}
        }
        delete(c.elem, key)
        c.lru.Remove(e)
      }
    }
    v = &cacheElement{key:key, value:value}
    v.elem = c.lru.PushFront(key)
    c.elem[key] = v
  }
}
