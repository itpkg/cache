package cache

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/itpkg/log"
)

type RedisStore struct {
	Pool   *redis.Pool `inject:""`
	Logger log.Logger  `inject:""`
	Prefix string      `inject:"cache prefix"`
}

func (p *RedisStore) Get(key string, val interface{}) error {
	c := p.Pool.Get()
	defer c.Close()
	if buf, err := redis.Bytes(c.Do("GET", p.Prefix+key)); err == nil {
		return json.Unmarshal(buf, val)
	} else {
		return err
	}

}

func (p *RedisStore) Set(key string, val interface{}, expire uint) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c := p.Pool.Get()
	defer c.Close()
	if expire == 0 {
		_, err = c.Do("SET", p.Prefix+key, buf)
	} else {
		_, err = c.Do("SET", p.Prefix+key, buf, "EX", expire)
	}

	return err
}

func (p *RedisStore) Delete(key string) error {
	c := p.Pool.Get()
	defer c.Close()
	_, err := c.Do("DEL", p.Prefix+key)
	return err
}

func (p *RedisStore) Flush() error {
	c := p.Pool.Get()
	defer c.Close()
	val, err := c.Do("KEYS", p.Prefix+"*")
	if err == nil {
		if ks := val.([]interface{}); len(ks) > 0 {
			_, err = c.Do("DEL", ks...)
		}
	}
	return err
}
