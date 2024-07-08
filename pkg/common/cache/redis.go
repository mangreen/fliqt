package cache

import (
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/charmbracelet/log"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

func InitCaches() (redis.Conn, error) {
	pool := &redis.Pool{
		MaxIdle:   4,
		MaxActive: 4,
		Dial: func() (redis.Conn, error) {
			rc, err := redis.Dial("tcp", viper.GetString("redis.address"))
			if err != nil {
				return nil, err
			}

			return rc, nil
		},
		IdleTimeout: time.Second,
		Wait:        true,
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	var conn redis.Conn

	backoff.Retry(func() error {
		conn = pool.Get()
		if conn == nil {
			log.Errorf("[redis] redis.Conn get failed")
			return errors.New("redis.Conn get failed")
		}

		ping, err := redis.String(conn.Do("ping"))
		if err != nil {
			log.Errorf("[redis] ping failed: %v", err)
			return err
		}

		log.Infof("[redis] redis.Conn PING success: %s", ping)

		return nil
	}, bo)

	if _, err := SelectDB(conn, 15); err != nil {
		return nil, err
	}

	return conn, nil
}

func SelectDB(conn redis.Conn, idx int) (redis.Conn, error) {
	_, err := conn.Do("SELECT", idx)
	if err != nil {
		log.Errorf("[redis] SELECT db %s failed: %v", idx, err)
		return nil, err
	}

	return conn, nil
}

func Get(conn redis.Conn, key string) (string, error) {
	res, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Errorf("[redis] GET %s failed: %v", key, err)
		return "", err
	}

	return res, nil
}

func Set(conn redis.Conn, key string, val interface{}) (string, error) {
	res, err := redis.String(conn.Do("SET", key, val))
	if err != nil {
		log.Errorf("[redis] SET %s failed: %v", key, err)
		return "", err
	}

	return res, nil
}

func HMGet(conn redis.Conn, key string, res interface{}) error {
	vals, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		log.Errorf("[redis] HGETALL %s failed: %v", key, err)
		return err
	}

	if err := redis.ScanStruct(vals, res); err != nil {
		log.Errorf("[redis] ScanStruct failed: %v", err)
		return err
	}

	return nil
}

func HMSet(conn redis.Conn, key string, hm interface{}) (string, error) {
	res, err := redis.String(conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(hm)...))
	if err != nil {
		log.Errorf("[redis] HMSET %s failed: %v", key, err)
		return "", err
	}

	return res, nil
}
