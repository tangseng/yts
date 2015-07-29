package models

import (
	"github.com/garyburd/redigo/redis"
	"crypto/md5"
	"fmt"
	"math/rand"
	"io"
	"time"
)

type RedisOption struct {
	num int
	size int
	life int64
	host string
	password string
	Pool *redis.Pool
}

func NewRedisOption(host string, password string, size int, num int, life int64) (*RedisOption, error){
	ro := &RedisOption{
		num : num,
		size : size,
		life : life,
		host : host,
		password : password,
	}
	ro.Pool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", ro.host)
		if err != nil {
			return nil, err
		}
		if ro.password != "" {
			if _, err := c.Do("AUTH", ro.password); err != nil {
				c.Close()
				return nil, err
			}
		}
		_, err = c.Do("SELECT", ro.num)
		if err != nil {
			c.Close()
			return nil, err
		}
		return c, err
	}, ro.size)
	err := ro.Pool.Get().Err()
	if err != nil {
		return nil, err
	}
	return ro, nil
}

var RO *RedisOption

func init() {
	var err error
	RO, err = NewRedisOption("127.0.0.1:6379", "vikings", 100, 0, 604800)
	if err != nil {
		panic("ttt")
	}
}

func MD5(str, salt string) string {
	h := md5.New()
	io.WriteString(h, str)
	io.WriteString(h, salt)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func RandStr(length int, replaceWords string) string {
	rand.Seed(Now())
	wordStr := "abcdefghijklmnopqrstuvwxyz"
	if len(replaceWords) != 0 {
		wordStr = replaceWords
	}
	words := []byte(wordStr)
	wLen := len(words)
	rs := make([]byte, 0)
	for i := 0; i < length; i++ {
		w := words[rand.Intn(wLen)]
		rs = append(rs, w)
	}
	return string(rs)
}

func RandStr2(length int) string {
	return RandStr(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

func Now() int64 {
	return time.Now().Unix()
}

func NowNano() int64 {
	return time.Now().UnixNano()
}