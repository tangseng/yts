package models

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"strconv"
)

type User struct {
	Name string
	Nick string
	Password string
	Salt string
	Time int64
}

func NewUser(name, nick, password, salt string, time int64) *User {
	if time == 0 {
		time = Now()
	}
	return &User{
		Name : name,
		Nick : nick,
		Password : password,
		Salt : salt,
		Time : time,
	}
}

func (this *User) CheckPass(password string) bool {
	return password == this.Password
}

type UserOption struct {
	userListKey string
	userKey string
	RO *RedisOption
}

func NewUserOption() *UserOption {
	return &UserOption{
		userListKey : "user:userlist",
		userKey : "user:user_",
		RO : RO,
	}
}

func (this *UserOption) Get(name string) (*User, error) {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey(name)
	if existed, _ := redis.Int(c.Do("EXISTS", key)); existed == 0 {
		return nil, fmt.Errorf("不存在")
	}
	userMap, err := redis.StringMap(c.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	time, _ :=strconv.ParseInt(userMap["time"], 10, 64)
	return NewUser(userMap["name"], userMap["nick"], userMap["password"], userMap["salt"], time), nil
}

func (this *UserOption) GetAll() []*User {
	c := this.RO.Pool.Get()
	defer c.Close()
	length, _ := redis.Int(c.Do("LLEN", this.userListKey))
	users := make([]*User, length)
	userNames, _ := redis.Strings(c.Do("LRANGE", this.userListKey, 0, length - 1))
	for k, userName := range userNames {
		userMap, _ := redis.StringMap(c.Do("HGETALL", this.nameKey(userName)))
		time, _ :=strconv.ParseInt(userMap["time"], 10, 64)
		users[k]= NewUser(userMap["name"], userMap["nick"], userMap["password"], userMap["salt"], time)
	}
	return users
}

func (this *UserOption) Create(user *User) error {
	c := this.RO.Pool.Get()
	defer c.Close()
	_, err := c.Do("HMSET", this.nameKey(user.Name),
		"name", user.Name,
		"nick", user.Nick,
		"password", user.Password,
		"salt", user.Salt,
		"time", strconv.FormatInt(user.Time, 10),
	)
	if err != nil {
		return err
	}
	_, err = c.Do("LPUSH", this.userListKey, user.Name)
	if err != nil {
		return err
	}
	return nil
}

func (this *UserOption) Update(user *User) error {
	c := this.RO.Pool.Get()
	defer c.Close()
	_, err := c.Do("HMSET", this.nameKey(user.Name),
		"nick", user.Nick,
		"password", user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (this *UserOption) Delete(name string) error {
	c := this.RO.Pool.Get()
	defer c.Close()
	outErr := fmt.Errorf("不存在该用户")
	num, err := redis.Int(c.Do("LREM", this.userListKey, 1, name))
	if num != 1 || err != nil {
		return outErr
	}
	num, err = redis.Int(c.Do("DEL", this.nameKey(name)))
	if num != 1 || err != nil {
		return outErr
	}
	return nil
}

func (this *UserOption) nameKey(name string) string {
	return fmt.Sprintf("%s%s", this.userKey, name)
}