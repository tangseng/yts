package models

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"github.com/astaxie/beego"
)

type Token struct {
	token string
	*User
}

func NewToken(user *User, token string) *Token {
	t := &Token{
		User : user,
	}
	if len(token) == 0 {
		token = MD5(fmt.Sprintf("%s%s%s", t.User.Name, t.User.Password, t.User.Salt), fmt.Sprintf("%d", Now()))
	}
	t.SetToken(token)
	return t
}

func (this *Token) SetToken(token string) {
	this.token = token
}

func (this *Token) Token() string {
	return this.token
}

func (this *Token) Val() string {
	return this.User.Name
}

type TokenOption struct {
	tokenHashKey string
	Expire int
	RO *RedisOption
}

func NewTokenOption() *TokenOption {
	expire, _ := beego.AppConfig.Int("Expire")
	return &TokenOption{
		tokenHashKey : "token:",
		Expire : expire,
		RO : RO,
	}
}

func (this *TokenOption) Get(token string) (*Token, error) {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey(token)
	if existed, _ := redis.Int(c.Do("EXISTS", key)); existed == 0 {
		return nil, fmt.Errorf("不存在")
	}
	userName, err := redis.String(c.Do("GET", key))
	if err != nil {
		return nil, err
	}
	user, _ := NewUserOption().Get(userName)
	t := NewToken(user, token)
	return t, nil
}

func (this *TokenOption) Add(t *Token) error {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey(t.Token())
	_, err := c.Do("SET", key, t.Val())
	if err != nil {
		return err
	}
	_, err = c.Do("EXPIRE", key, this.Expire)
	if err != nil {
		return err
	}
	return nil
}

func (this *TokenOption) Remove(token string) error {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey(token)
	_, err := c.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}

func (this *TokenOption) nameKey(name string) string {
	return fmt.Sprintf("%s%s", this.tokenHashKey, name)
}
