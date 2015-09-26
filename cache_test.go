package cache_test

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/itpkg/cache"
	"github.com/itpkg/log"
)

type Model struct {
	Message string
	Value   float64
	Now     time.Time
}

func TestRedis(t *testing.T) {
	var pool = &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	var store cache.Store
	store = &cache.RedisStore{
		Pool:   pool,
		Prefix: "cache://test/",
		Logger: log.NewStdoutLogger("test", log.DEBUG),
	}

	key := "test"
	val := &Model{Now: time.Now(), Value: 1.1, Message: "Hello"}

	if err := store.Set(key, val, 123); err != nil {
		t.Errorf("error on cache set: %v", err)
	}

	var tmp1 Model
	if err := store.Get(key, &tmp1); err == nil {
		t.Logf("%v => %v", val, tmp1)
		if tmp1.Message != val.Message {
			t.Errorf("bad data on set")
		}
	} else {
		t.Errorf("error on cache set: %v", err)
	}

	if err := store.Delete(key); err != nil {
		t.Errorf("error on cache del: %v", err)
	}

	if err := store.Set(key, val, 0); err != nil {
		t.Errorf("error on cache set: %v", err)
	}
	if err := store.Flush(); err != nil {
		t.Errorf("error on cache del: %v", err)
	}

}
