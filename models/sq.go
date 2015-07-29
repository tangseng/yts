package models

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"encoding/json"
)

type SQ struct {
	Name string `json:name`
	Nick string `json:nick`
	Password string `json:password`
	Status int8 `json:status`
}

func NewSQ(name, nick, password string) *SQ {
	return &SQ{
		Name : name,
		Nick : nick,
		Password : password,
	}
}

type SQOption struct {
	sqHashKey string
	RO *RedisOption
}

func NewSQOption() *SQOption {
	return &SQOption{
		sqHashKey : "sq:sqhash",
		RO : RO,
	}
}

func (this *SQOption) Get(name string) (*SQ, error) {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey()
	if existed, _ := redis.Int(c.Do("EXISTS", key)); existed == 0 {
		return nil, fmt.Errorf("no")
	}
	sqJsonString, err := redis.String(c.Do("HGET", key, name))
	if err != nil {
		return nil, err
	}
	sq := &SQ{}
	json.Unmarshal([]byte(sqJsonString), sq)
	return sq, nil
}

func (this *SQOption) GetAll() []*SQ {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey()
	sqSlice := make([]*SQ, 0)
	if existed, _ := redis.Int(c.Do("EXISTS", key)); existed == 0 {
		return sqSlice
	}
	sqMap, err := redis.StringMap(c.Do("HGETALL", key))
	if err != nil {
		return sqSlice
	}
	for _, sqJsonString := range sqMap {
		sq := SQ{}
		err = json.Unmarshal([]byte(sqJsonString), &sq)
		if err != nil {
			continue
		}
		sqSlice = append(sqSlice, &sq)
	}
	return sqSlice
}

func (this *SQOption) Set(sq *SQ) error {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey()
	sqBytes, err := json.Marshal(sq)
	if err != nil {
		return err
	}
	_, err = c.Do("HSET", key, sq.Name, string(sqBytes))
	if err != nil {
		return err
	}
	return nil
}

func (this *SQOption) Check(sq *SQ) bool {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey()
	if existed, _ := redis.Int(c.Do("HEXISTS", key, sq.Name)); existed == 0 {
		return true
	}
	return false
}

func (this *SQOption) nameKey() string {
	return fmt.Sprintf("%s", this.sqHashKey)
}
