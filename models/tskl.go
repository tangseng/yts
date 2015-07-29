package models

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

type TSKL struct {
	Send string
	Back string
	Berr string
	Hash string
	*User
}

func NewTSKL(user *User) *TSKL {
	return &TSKL{
		Hash : MD5(user.Name, strconv.FormatInt(Now(), 10)),
		User : user,
	}
}

type TSKLOption struct {
	tsklHashKey string
	Expire int
	RO *RedisOption
}

func NewTSKLOption() *TSKLOption {
	expire, _ := beego.AppConfig.Int("Expire")
	return &TSKLOption{
		tsklHashKey : "tskl:",
		Expire : expire,
		RO : RO,
	}
}

func (this *TSKLOption) Get(user *User) (*TSKL, error) {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey(user.Name)
	if existed, _ := redis.Int(c.Do("EXISTS", key)); existed == 0 {
		return nil, fmt.Errorf("不存在")
	}
	tsklMap, err := redis.StringMap(c.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	tskl := NewTSKL(user)
	tskl.Hash = tsklMap["hash"]
	tskl.Send = tsklMap["send"]
	tskl.Back = tsklMap["back"]
	tskl.Berr = tsklMap["berr"]
	return tskl, nil
}

func (this *TSKLOption) Set(tskl *TSKL) error {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey(tskl.User.Name)
	_, err := c.Do("HMSET", key, "hash", tskl.Hash, "send", tskl.Send, "back", tskl.Back, "berr", tskl.Berr)
	if err != nil {
		return err
	}
	_, err = c.Do("EXPIRE", key, this.Expire)
	if err != nil {
		return err
	}
	return nil
}

func (this *TSKLOption) SetWhich(tskl *TSKL, which string) error {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey(tskl.User.Name)
	var val string
	switch which {
	case "send":
		val = tskl.Send
	case "back":
		val = tskl.Back
	case "berr":
		val = tskl.Berr
	}
	_, err := c.Do("HMSET", key, which, val)
	if err != nil {
		return err
	}
	return nil
}

func (this *TSKLOption) nameKey(name string) string {
	return fmt.Sprintf("%s%s", this.tsklHashKey, name)
}
