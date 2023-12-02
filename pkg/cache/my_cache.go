package cache

import (
	message "awesomeProject"
	"github.com/patrickmn/go-cache"
	"log"
	"strconv"
	"time"
)

type myCache struct {
	messCache *cache.Cache
}

const (
	defaultExpiration = 5 * time.Minute
	purgeTime         = 10 * time.Minute
)

func newMyCache() *myCache {
	Cache := cache.New(defaultExpiration, purgeTime)
	return &myCache{messCache: Cache}
}

func (c *myCache) Read(id int) (item message.Message, ok bool) {
	product, ok := c.messCache.Get(strconv.Itoa(id))
	if ok {
		log.Println("from cache")
		res := product.(message.Message)
		return res, true
	}
	return message.Message{}, false
}
func (c *myCache) Update(id int, messageCache message.Message) {
	c.messCache.Set(strconv.Itoa(id), messageCache, cache.DefaultExpiration)
}

var C = newMyCache()
