package easycompact

import (
	"sync"
	"time"
)

type EasyCompact interface {
	Set(key any, data any)
	Close()
}

type l struct {
	data []any
	ttl  time.Time
	sync.RWMutex
}

type cache struct {
	list   map[any]*l
	update func(in any, data any)
	sync.RWMutex
	fn func(key any, data []any)
}

func New(ttl *time.Duration, fn func(key any, data []any)) EasyCompact {

	var c = cache{}
	c.list = make(map[any]*l)
	c.fn = fn

	go c.pushJob(ttl)

	return &c
}
func (c *cache) pushJob(ttl *time.Duration) {

	ticker := time.NewTicker(*ttl)
	for range ticker.C {
		c.push()
	}
}

func (c *cache) push() {
	list := make(map[any]*l)
	c.Lock()
	for k, v := range c.list {
		list[k] = v
	}
	c.list = make(map[any]*l)
	c.Unlock()

	for k, v := range c.list {
		k := k
		data := v.data

		delete(c.list, k)

		c.fn(k, data)

	}
}

func (c *cache) Set(key any, data any) {
	c.Lock()
	defer c.Unlock()

	d, ok := c.list[key]
	if !ok {
		c.list[key] = &l{}
		d = c.list[key]
	}

	d.data = append(d.data, data)
}

func (c *cache) Close() {
	c.push()
}
