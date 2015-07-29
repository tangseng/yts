package models

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"github.com/astaxie/beego"
	"encoding/json"
)

type TS struct {
	Data string `json:"data"`
	Ip string `json:"ip"`
	TSTime int64 `json:"tstime"`
	Time int64 `json:"time"`
	Type string `json:"type"`
}

func NewTS(data, t, ip string, tstime, time int64) *TS {
	if time == 0 {
		time = NowNano()
	}
	return &TS{
		Data : data,
		Ip : ip,
		TSTime : tstime,
		Time : time,
		Type : t,
	}
}

type TSOption struct {
	tsListKey string
	Limit int
	RO *RedisOption
}

func NewTSOption() *TSOption {
	limit, _ := beego.AppConfig.Int("Limit")
	return &TSOption{
		tsListKey : "ts:tslist_",
		Limit : limit,
		RO : RO,
	}
}

func (this *TSOption) Get(name string, start, num int) ([]TS, error) {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey(name)
	if existed, _ := redis.Int(c.Do("EXISTS", key)); existed == 0 {
		return nil, fmt.Errorf("不存在")
	}
	tss, err := redis.Strings(c.Do("LRANGE", key, start, start + num))
	if err != nil {
		return nil, err
	}
	outTss := make([]TS, 0)
	for _, tsv := range tss {
		ts := TS{}
		json.Unmarshal([]byte(tsv), &ts)
		outTss = append(outTss, ts)
	}
	return outTss, nil
}

func (this *TSOption) Insert(name string, ts *TS) error {
	c := this.RO.Pool.Get()
	defer c.Close()
	key := this.nameKey(name)
	jsonBS, _ := json.Marshal(*ts)
	_, err := c.Do("LPUSH", key, string(jsonBS))
	if err != nil {
		return err
	}
	_, err = c.Do("LTRIM", key, 0, this.Limit)
	if err != nil {
		return err
	}
	return nil
}

func (this *TSOption) nameKey(name string) string {
	return fmt.Sprintf("%s%s", this.tsListKey, name)
}
