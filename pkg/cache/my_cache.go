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
	//Продолжительность жизни кеша по-умолчанию
	defaultExpiration = 5 * time.Minute
	//Интервал, через который запускается механизм очистки кеша
	purgeTime = 10 * time.Minute
)

func newMyCache() *myCache {
	Cache := cache.New(defaultExpiration, purgeTime)
	return &myCache{messCache: Cache}
}

// Функция для получения объекта Message из кэша по id
func (c *myCache) Read(id int) (item message.Message, ok bool) {
	product, ok := c.messCache.Get(strconv.Itoa(id))
	if ok {
		log.Println("from cache")
		res := product.(message.Message)
		return res, true
	}
	return message.Message{}, false
}

// Функция для сохранения объекта Message в кэш
func (c *myCache) Update(id int, messageCache message.Message) {
	c.messCache.Set(strconv.Itoa(id), messageCache, cache.DefaultExpiration)
}

// Переменная хранящая экземпляр myCach
var C = newMyCache()
