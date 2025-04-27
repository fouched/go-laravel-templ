package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Cache interface {
	Has(string) (bool, error)
	Get(string) (interface{}, error)
	Set(string, interface{}, ...int) error
	Forget(string) error
	EmptyByMatch(string) error
	Empty() error
}

type RedisCache struct {
	Conn   *redis.Client
	Prefix string
}

type Entry map[string]interface{}

func (c *RedisCache) Has(str string) (bool, error) {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Conn

	val, err := conn.Get(ctx, key).Result()

	switch {
	case errors.Is(err, redis.Nil):
		fmt.Println("key does not exist")
		return false, err
	case err != nil:
		fmt.Println("get failed", err)
		return false, err
	case val == "":
		fmt.Println("value is empty")
		return false, err
	default:
		return true, nil
	}
}

func (c *RedisCache) Get(str string) (interface{}, error) {
	return "", nil
}

func (c *RedisCache) Set(str string, value interface{}, ttl ...int) error {

	return nil
}

func (c *RedisCache) Forget(str string) error {

	return nil
}

func (c *RedisCache) EmptyByMatch(str string) error {

	return nil
}

func (c *RedisCache) Empty() error {

	return nil
}
