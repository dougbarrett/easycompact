package easycompact

import (
	"sync"
	"time"
)

func New(ttl *time.Duration, fn func(key any, data []any)) func(key any, data any) {
	type l struct {
		data []any
		ttl  time.Time
		sync.RWMutex
	}

	type cache struct {
		list   map[any]*l
		update func(in any, data any)
		sync.RWMutex
	}

	var c = cache{}
	c.list = make(map[any]*l)

	ticker := time.NewTicker(*ttl)

	go func(c *cache) {
		for range ticker.C {
			for k, v := range c.list {
				k := k
				data := v.data

				delete(c.list, k)

				fn(k, data)

			}
		}
	}(&c)

	return func(key any, data any) {
		c.Lock()
		defer c.Unlock()

		d, ok := c.list[key]
		if !ok {
			c.list[key] = &l{}
			d = c.list[key]
		}

		d.data = append(d.data, data)
	}
}
