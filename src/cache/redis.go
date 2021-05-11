package cache

import (
	"awesomeProject/src/config"
	red "github.com/gomodule/redigo/redis"
	"time"
)

type Redis struct {
	pool *red.Pool
}

var rdc *Redis

func InitRedis(redisConfig config.RedisConfig) error {
	rdc = new(Redis)
	rdc.pool = &red.Pool{
		MaxIdle:     redisConfig.MaxIdle,
		MaxActive:   redisConfig.MaxActive,
		IdleTimeout: time.Duration(redisConfig.IdleTimeout),
		Dial: func() (red.Conn, error) {
			return red.Dial(
				redisConfig.Network,
				redisConfig.Address,
				red.DialReadTimeout(time.Duration(redisConfig.DialReadTimeout)*time.Microsecond),
				red.DialWriteTimeout(time.Duration(redisConfig.DialWriteTimeout)*time.Microsecond),
				red.DialConnectTimeout(time.Duration(redisConfig.DialConnectTimeout)*time.Microsecond),
				red.DialDatabase(redisConfig.Database),
			)
		},
	}
	return nil
}

func Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	conn := rdc.pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	defer conn.Close()
	params := make([]interface{}, 0)
	params = append(params, key)
	if len(args) > 0 {
		for _, v := range args {
			params = append(params, v)
		}
	}
	return conn.Do(cmd, params...)
}
