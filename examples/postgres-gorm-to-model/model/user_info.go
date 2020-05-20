package model

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// Auto-Generate Header
/*
    Code is auto-generated by github.com/fwhezfwhez/model_convert.Package below might be imported:
      - github.com/fwhezfwhez/errorx
      - github.com/garyburd/redigo/redis
      - github.com/jinzhu/gorm
    You can get them by:
      - go get github.com/fwhezfwhez/errorx
      - go get github.com/garyburd/redigo/redis
      - go get github.com/jinzhu/gorm

    To fulfill redis part, don't forget to set TODOs.They are:
      - RedisKey() string
      - RedisSecondDuration() int
*/
type UserInfo struct {
	Id            int       `gorm:"column:id;default:" json:"id" form:"id"`
	UserId        int       `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
	OpenId        string    `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`
	UnionId       string    `gorm:"column:union_id;default:" json:"union_id" form:"union_id"`
	UserName      string    `gorm:"column:user_name;default:" json:"user_name" form:"user_name"`
	HeaderUrl     string    `gorm:"column:header_url;default:" json:"header_url" form:"header_url"`
	Sex           int       `gorm:"column:sex;default:" json:"sex" form:"sex"`
	GameId        int       `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
	LastLoginTime time.Time `gorm:"column:last_login_time;default:" json:"last_login_time" form:"last_login_time"`
	CreatedAt     time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
	PlatformId    int       `gorm:"column:platform_id;default:" json:"platform_id" form:"platform_id"`
	Channel       string    `gorm:"column:channel;default:" json:"channel" form:"channel"`
	IsRobot       int       `gorm:"column:is_robot;default:" json:"is_robot" form:"is_robot"`
}

func (o UserInfo) TableName() string {
	return "user_info"
}

var UserInfoRedisKeyFormat = "xyx:user_info:%s:%d:%d"

func (o UserInfo) RedisKey() string {
	// TODO set its redis key and required args
	return fmt.Sprintf(UserInfoRedisKeyFormat, "pro", o.GameId, o.UserId)
}

var ArrayUserInfoRedisKeyFormat = "xyx:user_infos:%s:%d"

func (o UserInfo) ArrayRedisKey() string {
	// TODO set its array key and required args
	return fmt.Sprintf(ArrayUserInfoRedisKeyFormat, "pro", o.GameId)
}

func (o UserInfo) RedisSecondDuration() int {
	// TODO set its redis duration, default 1-7 day,  return -1 means no time limit
	return int(time.Now().Unix()%7+1) * 24 * 60 * 60
}

// TODO,set using db or not. If set false, o.MustGet() will never get its data from db.
func (o UserInfo) UseDB() bool {
	return false
}

func (o *UserInfo) GetFromRedis(conn redis.Conn) error {
	if o.RedisKey() == "" {
		return errorx.NewFromString("object UserInfo has not set redis key yet")
	}
	buf, e := redis.Bytes(conn.Do("GET", o.RedisKey()))

	if e == nil && string(buf) == "DISABLE" {
		return fmt.Errorf("not found record in db nor redis")
	}

	if e == redis.ErrNil {
		return e
	}

	if e != nil && e != redis.ErrNil {
		return errorx.Wrap(e)
	}

	e = json.Unmarshal(buf, &o)

	if e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

func (o *UserInfo) ArrayGetFromRedis(conn redis.Conn) ([]UserInfo, error) {
	if o.ArrayRedisKey() == "" {
		return nil, errorx.NewFromString("object UserInfo has not set redis key yet")
	}

	var list = make([]UserInfo, 0, 10)
	buf, e := redis.Bytes(conn.Do("GET", o.ArrayRedisKey()))

	// avoid passing through and hit database
	// When o.ArrayMustGet() not found both in redis and db, will set its key DISABLE
	// and return 'fmt.Errorf("not found record in db nor redis")'
	if e == nil && string(buf) == "DISABLE" {
		return nil, fmt.Errorf("not found record in db nor redis")
	}

	// Not found in redis
	if e == redis.ErrNil {
		return nil, e
	}

	// Server error, should be logged by caller
	if e != nil && e != redis.ErrNil {
		return nil, errorx.Wrap(e)
	}

	e = json.Unmarshal(buf, &list)

	if e != nil {
		return nil, errorx.Wrap(e)
	}
	return list, nil
}

// engine should prepare its condition.
// if record not found,it will return 'var notFound = fmt.Errorf("not found record in db nor redis")'.
// If you want to ignore not found error, do it like:
// if e:= o.MustGet(conn, engine.Model(Model{}).Where("condition =?", arg)).Error;e!=nil {
//     if e.Error() == "not found record in db nor redis"{
//         log.Println(e)
//         return
//     }
// }
func (o *UserInfo) MustGet(conn redis.Conn, engine *gorm.DB) error {
	var shouldSyncToCache bool

	if UserInfoCacheSwitch {
		if e := o.GetFromCache(); e == nil {
			return nil
		}
		defer func() {
			if shouldSyncToCache {
				fmt.Println("exec sync to cache")
				o.SyncToCache()
			}
		}()
	}

	e := o.GetFromRedis(conn)
	// When redis key stores its value 'DISABLE', will returns notFoundError and no need to query from db any more
	if e != nil && e.Error() == "not found record in db nor redis" {
		return e
	}

	if e == nil {
		shouldSyncToCache = true
		return nil
	}
	if e != nil {
		var count int
		if e2 := engine.Count(&count).Error; e2 != nil {
			return errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e2))
		}
		if count == 0 {
			var notFound = fmt.Errorf("not found record in db nor redis")
			if o.RedisSecondDuration() == -1 {
				conn.Do("SET", o.RedisKey(), "DISABLE", "NX")
			} else {
				conn.Do("SET", o.RedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			}
			return notFound
		}

		if e3 := engine.First(&o).Error; e3 != nil {
			return errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e3))
		}
		shouldSyncToCache = true

		if e == redis.ErrNil {
			o.SyncToRedis(conn)
			return nil
		}
		return errorx.Wrap(e)
	}
	return nil
}

func (o *UserInfo) ArrayMustGet(conn redis.Conn, engine *gorm.DB) ([]UserInfo, error) {
	var shouldSyncToCache bool
	var arr []UserInfo

	if ArrayUserInfoCacheSwitch {
		if arr, e := o.ArrayGetFromCache(); e == nil {
			return arr, nil
		}
		defer func() {
			if shouldSyncToCache {
				fmt.Println("exec sync to cache")
				o.ArraySyncToCache(arr)
			}
		}()
	}

	list, e := o.ArrayGetFromRedis(conn)
	// When redis key stores its value 'DISABLE', will returns notFoundError and no need to query from db any more
	// When call ArrayDeleteFromRedis(), will activate its redis and db query
	if e != nil && e.Error() == "not found record in db nor redis" {
		return nil, e
	}
	// get from redis success.
	if e == nil {
		shouldSyncToCache = true
		arr = list
		return list, nil
	}
	// get from redis fail, try db
	if e != nil {
		var count int
		if e2 := engine.Count(&count).Error; e2 != nil {
			return nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e2))
		}
		if count == 0 {
			var notFound = fmt.Errorf("not found record in db nor redis")
			if o.RedisSecondDuration() == -1 {
				conn.Do("SET", o.ArrayRedisKey(), "DISABLE", "NX")
			} else {
				conn.Do("SET", o.ArrayRedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			}
			return nil, notFound
		}

		if e3 := engine.Find(&list).Error; e3 != nil {
			return nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e3))
		}

		shouldSyncToCache = true
		arr = list
		// try sync to redis
		if e == redis.ErrNil {
			o.ArraySyncToRedis(conn, list)
			return list, nil
		}
		return nil, errorx.Wrap(e)
	}
	return nil, nil
}

func (o UserInfo) SyncToRedis(conn redis.Conn) error {
	if o.RedisKey() == "" {
		return errorx.NewFromString("object UserInfo has not set redis key yet")
	}
	buf, e := json.Marshal(o)
	if e != nil {
		return errorx.Wrap(e)
	}
	if o.RedisSecondDuration() == -1 {
		if _, e := conn.Do("SET", o.RedisKey(), buf); e != nil {
			return errorx.Wrap(e)
		}
	} else {
		if _, e := conn.Do("SETEX", o.RedisKey(), o.RedisSecondDuration(), buf); e != nil {
			return errorx.Wrap(e)
		}
	}
	return nil
}

func (o UserInfo) ArraySyncToRedis(conn redis.Conn, list []UserInfo) error {
	if o.ArrayRedisKey() == "" {
		return errorx.NewFromString("object UserInfo has not set redis key yet")
	}
	buf, e := json.Marshal(list)
	if e != nil {
		return errorx.Wrap(e)
	}
	if o.RedisSecondDuration() == -1 {
		if _, e := conn.Do("SET", o.ArrayRedisKey(), buf); e != nil {
			return errorx.Wrap(e)
		}
	} else {
		if _, e := conn.Do("SETEX", o.ArrayRedisKey(), o.RedisSecondDuration(), buf); e != nil {
			return errorx.Wrap(e)
		}
	}
	return nil
}

func (o UserInfo) DeleteFromRedis(conn redis.Conn) error {
	if o.RedisKey() != "" {
		if _, e := conn.Do("DEL", o.RedisKey()); e != nil {
			return errorx.Wrap(e)
		}
	}

	if o.ArrayRedisKey() != "" {
		if _, e := conn.Do("DEL", o.ArrayRedisKey()); e != nil {
			return errorx.Wrap(e)
		}
	}

	if UserInfoCacheSwitch {
		o.deleteFromCache()
	}
	if ArrayUserInfoCacheSwitch {
		o.ArraydeleteFromCache()
	}
	return nil
}
func (o UserInfo) ArrayDeleteFromRedis(conn redis.Conn) error {
	return o.DeleteFromRedis(conn)
}

// Dump data through api GET remote url generated by 'GenerateListApi()' to local database.
// This method should never used in production. It's best to to run it before app is running.
//
// mode=1, each time will delete old local data and dump from api.
// mode=2, each time will update/keep the existed data. Mode=2 is developing.
func (o UserInfo) DumpToLocal(url string, engine *gorm.DB, mode int) error {
	tableName := o.TableName()

	tran := engine.Begin()
	if e := tran.Exec(fmt.Sprintf("delete from %s", tableName)).Error; e != nil {
		tran.Rollback()
		return errorx.Wrap(e)
	}

	type Result struct {
		Data  []UserInfo `json:"data"`
		Count int        `json:"count"`
	}
	var result Result
	resp, e := http.Get(url)
	if e != nil {
		tran.Rollback()
		return errorx.Wrap(e)
	}
	if resp == nil || resp.Body == nil {
		tran.Rollback()
		return errorx.NewFromString("resp or body nil")
	}
	defer resp.Body.Close()

	buf, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		tran.Rollback()
		return errorx.Wrap(e)
	}

	if resp.StatusCode != 200 {
		var body string
		if len(buf) < 100 {
			body = string(buf)
		} else {
			body = string(buf[:100])
		}
		return errorx.NewFromStringf("status not 200, got %d,body %s", resp.StatusCode, body)
	}

	if e := json.Unmarshal(buf, &result); e != nil {
		tran.Rollback()
		return errorx.Wrap(e)
	}

	for i, _ := range result.Data {
		data := result.Data[i]
		if e := tran.Model(&o).Create(&data).Error; e != nil {
			tran.Rollback()
			return errorx.Wrap(e)
		}
	}
	tran.Commit()
	return nil
}

// 2nd-cache Header
// 2nd-cache share RedisKey() as its key.

// self Header
var (
	UserInfoCache         = make(map[string]UserInfo, 0)
	UserInfoCacheKeyOrder = make([]string, 0, 500)

	UserInfoCacheLock = sync.RWMutex{}

	UserInfoNotFoundErr  = fmt.Errorf("not found in cache")
	UserInfoSwitchOffErr = fmt.Errorf("2nd-cache switch is off")
)

const (
	// Max cap of map and len
	UserInfoCacheMaxLength = 5000
	// When faced with max cap, will delete 2000/5 data from map and arr
	// DeleteRate can't be zero.
	UserInfoDeleteRate = 5
	// Whether use cache
	UserInfoCacheSwitch = false
)

func (o *UserInfo) GetFromCache() error {
	if UserInfoCacheSwitch == false {
		return UserInfoSwitchOffErr
	}
	UserInfoCacheLock.RLock()
	defer UserInfoCacheLock.RUnlock()
	tmp, ok := UserInfoCache[o.RedisKey()]
	if !ok {
		return UserInfoNotFoundErr
	}
	*o = tmp
	fmt.Println("get from cache")
	return nil
}

func (o *UserInfo) deleteFromCache() {
	if UserInfoCacheSwitch == false {
		return
	}
	UserInfoCacheLock.Lock()
	defer UserInfoCacheLock.Unlock()

	delete(UserInfoCache, o.RedisKey())
}

func (o *UserInfo) SyncToCache() {
	if UserInfoCacheSwitch == false {
		return
	}

	if UserInfoDeleteRate == 0 || UserInfoDeleteRate < 0 {
		return
	}

	if UserInfoCacheMaxLength == 0 {
		return
	}

	UserInfoCacheLock.Lock()
	defer UserInfoCacheLock.Unlock()

	leng := len(UserInfoCacheKeyOrder)
	if leng >= UserInfoCacheMaxLength {
		delta := UserInfoCacheMaxLength / UserInfoDeleteRate
		for i := 0; i < delta; i++ {
			if _, ok := UserInfoCache[UserInfoCacheKeyOrder[i]]; ok {
				delete(UserInfoCache, UserInfoCacheKeyOrder[i])
			}
		}
		UserInfoCacheKeyOrder = UserInfoCacheKeyOrder[delta:]
	}
	UserInfoCache[o.RedisKey()] = *o
	UserInfoCacheKeyOrder = append(UserInfoCacheKeyOrder, o.RedisKey())
}

// self Tail

// array Header
var (
	ArrayUserInfoCache         = make(map[string][]UserInfo, 0)
	ArrayUserInfoCacheKeyOrder = make([]string, 0, 500)

	ArrayUserInfoCacheLock = sync.RWMutex{}

	ArrayUserInfoNotFoundErr  = fmt.Errorf("not found in cache")
	ArrayUserInfoSwitchOffErr = fmt.Errorf("2nd-cache switch is off")
)

const (
	// Max cap of map and len
	ArrayUserInfoCacheMaxLength = 5000
	// When faced with max cap, will delete 2000/5 data from map and arr
	// DeleteRate can't be zero.
	ArrayUserInfoDeleteRate = 5
	// Whether use cache
	ArrayUserInfoCacheSwitch = false
)

func (o *UserInfo) ArrayGetFromCache() ([]UserInfo, error) {
	if ArrayUserInfoCacheSwitch == false {
		return nil, ArrayUserInfoSwitchOffErr
	}
	ArrayUserInfoCacheLock.RLock()
	defer ArrayUserInfoCacheLock.RUnlock()
	tmp, ok := ArrayUserInfoCache[o.ArrayRedisKey()]
	if !ok {
		return nil, ArrayUserInfoNotFoundErr
	}
	fmt.Println("get from cache")
	return tmp, nil
}

func (o *UserInfo) ArraydeleteFromCache() {
	if ArrayUserInfoCacheSwitch == false {
		return
	}
	ArrayUserInfoCacheLock.Lock()
	defer ArrayUserInfoCacheLock.Unlock()

	delete(ArrayUserInfoCache, o.ArrayRedisKey())
}

func (o *UserInfo) ArraySyncToCache(arr []UserInfo) {
	if ArrayUserInfoCacheSwitch == false {
		return
	}

	if ArrayUserInfoDeleteRate == 0 || ArrayUserInfoDeleteRate < 0 {
		return
	}

	if ArrayUserInfoCacheMaxLength == 0 {
		return
	}

	ArrayUserInfoCacheLock.Lock()
	defer ArrayUserInfoCacheLock.Unlock()

	leng := len(ArrayUserInfoCacheKeyOrder)
	if leng >= ArrayUserInfoCacheMaxLength {
		delta := ArrayUserInfoCacheMaxLength / ArrayUserInfoDeleteRate
		for i := 0; i < delta; i++ {
			if _, ok := ArrayUserInfoCache[ArrayUserInfoCacheKeyOrder[i]]; ok {
				delete(ArrayUserInfoCache, ArrayUserInfoCacheKeyOrder[i])
			}
		}
		ArrayUserInfoCacheKeyOrder = ArrayUserInfoCacheKeyOrder[delta:]
	}
	ArrayUserInfoCache[o.ArrayRedisKey()] = arr
	ArrayUserInfoCacheKeyOrder = append(ArrayUserInfoCacheKeyOrder, o.ArrayRedisKey())
}

// array Tail

// 2nd-cache Tail

// flexible-cache Header
// func (o UserInfo) ${cache_name}Key() string{
// 	// TODO-Set cache redis key
// 	return ""
// }
// func (o UserInfo) ${cache_name}Duration() int{
// 	// TODO-Set cache redis key expire duration. Default 1-7 days
//     return int(time.Now().Unix() % 7 + 1) * 24 * 60 * 60
// }
// func (o *UserInfo) ${cache_name}MustGet(conn redis.Conn, source func(${cache_name} *${cache_type})error) (${cache_type}, error) {

// 	rs, e:= redis.${Cache_type}(conn.Do("GET", o.${cache_name}Key()))
// 	if e !=nil {
// 		if e == redis.ErrNil {
//             if e:=source(&rs); e!=nil {
// 				return rs, errorx.Wrap(e)
// 			}
// 			if _, e= conn.Do("SETEX",  o.${cache_name}Key(), ${cache_name}Duration(), rs),; e!=nil {
// 				return rs, errorx.Wrap(e)
// 			}
// 			return rs,nil
// 		}
// 		return rs, errorx.Wrap(e)
// 	}
// 	return rs,nil

// }
// flexible-cache Tail

// no-decode Header
//
// MustGetNoDecode do most similar work as MustGet do, but it will not unmarshal data from redis into 'o', in the meanwhile, will return its raw json stream as return.
// This function aims to save cost of decoding in the only case that you want to return 'o' itself and has nothing changed to inner values.
// 'engine' should prepare its condition.
// if record not found,it will return 'var notFound = fmt.Errorf("not found record in db nor redis")'.
// If you want to ignore not found error, do it like:
// if buf, e:= o.MustGetNoDecode(conn, engine.Model(Model{}).Where("condition =?", arg)).Error;e!=nil {
//     if e.Error() == "not found record in db nor redis" || e == redis.ErrNil {
//         log.Println(e)
//         return
//     }
// }
//
func (o *UserInfo) MustGetNoDecode(conn redis.Conn, engine *gorm.DB) (json.RawMessage, error) {
	var shouldSyncToCache bool

	if UserInfoCacheSwitch {
		if e := o.GetFromCache(); e == nil {
			return nil, nil
		}
		defer func() {
			if shouldSyncToCache {
				fmt.Println("exec sync to cache")
				o.SyncToCache()
			}
		}()
	}

	arrBuf, e := o.GetFromRedisNoDecode(conn)
	// When redis key stores its value 'DISABLE', will returns notFoundError and no need to query from db any more
	if e != nil && e.Error() == "not found record in db nor redis" {
		return nil, e
	}

	if e == nil {
		shouldSyncToCache = true
		return arrBuf, nil
	}
	if e != nil {
		var count int
		if e2 := engine.Count(&count).Error; e2 != nil {
			return nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e2))
		}
		if count == 0 {
			var notFound = fmt.Errorf("not found record in db nor redis")
			if o.RedisSecondDuration() == -1 {
				conn.Do("SET", o.RedisKey(), "DISABLE", "NX")
			} else {
				conn.Do("SET", o.RedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			}
			return nil, notFound
		}

		if e3 := engine.First(&o).Error; e3 != nil {
			return nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e3))
		}
		shouldSyncToCache = true

		if e == redis.ErrNil {
			o.SyncToRedis(conn)
			return nil, nil
		}
		return nil, errorx.Wrap(e)
	}
	return nil, nil
}

// GetFromRedisNoDecode will return its json raw stream and will not decode into 'o'.
// It aims to save cost of decoding if json stream is decoded slowly.
func (o *UserInfo) GetFromRedisNoDecode(conn redis.Conn) (json.RawMessage, error) {
	if o.RedisKey() == "" {
		return nil, errorx.NewFromString("object UserInfo has not set redis key yet")
	}
	buf, e := redis.Bytes(conn.Do("GET", o.RedisKey()))

	if e == nil && string(buf) == "DISABLE" {
		return nil, fmt.Errorf("not found record in db nor redis")
	}

	if e == redis.ErrNil {
		return nil, e
	}

	if e != nil && e != redis.ErrNil {
		return nil, errorx.Wrap(e)
	}

	return buf, nil
}

// ArrayMustGetNoDecode will not unmarshal json stream to 'arr' and return json.Rawmessage as return value instead if it's found in redis,
// otherwise will return arr from cache or db.
//
// This function aims to save cost of decoding in the read-only case of 'o'. It means you should do nothing changed to its json value.
/*
	arr, arrBuf, e:= o.ArrayMustGetNoDecode(conn, engine)
	if e!=nil {
	// handle error
	}

	if len(arrBuf) >0 {
	c.JSON(200, gin.H{"message":"success", "data": arrBuf})
	} else {
		c.JSON(200, gin.H{"message":"success", "data": arr})
	}
*/
func (o *UserInfo) ArrayMustGetNoDecode(conn redis.Conn, engine *gorm.DB) ([]UserInfo, json.RawMessage, error) {
	var shouldSyncToCache bool
	var arr []UserInfo

	if ArrayUserInfoCacheSwitch {
		if arr, e := o.ArrayGetFromCache(); e == nil {
			return arr, nil, nil
		}
		defer func() {
			if shouldSyncToCache {
				fmt.Println("exec sync to cache")
				o.ArraySyncToCache(arr)
			}
		}()
	}

	arrBuf, e := o.ArrayGetFromRedisNoDecode(conn)
	// When redis key stores its value 'DISABLE', will returns notFoundError and no need to query from db any more
	// When call ArrayDeleteFromRedis(), will activate its redis and db query
	if e != nil && e.Error() == "not found record in db nor redis" {
		return nil, nil, e
	}
	// get from redis success.
	if e == nil {
		// shouldSyncToCache = true
		// arr = list
		return nil, arrBuf, nil
	}
	// get from redis fail, try db
	if e != nil {
		var list = make([]UserInfo, 0, 100)
		var count int
		if e2 := engine.Count(&count).Error; e2 != nil {
			return nil, nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e2))
		}
		if count == 0 {
			var notFound = fmt.Errorf("not found record in db nor redis")
			if o.RedisSecondDuration() == -1 {
				conn.Do("SET", o.ArrayRedisKey(), "DISABLE", "NX")
			} else {
				conn.Do("SET", o.ArrayRedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			}
			return nil, nil, notFound
		}

		if e3 := engine.Find(&list).Error; e3 != nil {
			return nil, nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e3))
		}

		shouldSyncToCache = true
		arr = list
		// try sync to redis
		if e == redis.ErrNil {
			o.ArraySyncToRedis(conn, list)
			return list, nil, nil
		}
		return nil, nil, errorx.Wrap(e)
	}
	return nil, nil, nil
}

func (o *UserInfo) ArrayGetFromRedisNoDecode(conn redis.Conn) (json.RawMessage, error) {
	if o.ArrayRedisKey() == "" {
		return nil, errorx.NewFromString("object UserInfo has not set redis key yet")
	}

	buf, e := redis.Bytes(conn.Do("GET", o.ArrayRedisKey()))

	// avoid passing through and hit database
	// When o.ArrayMustGet() not found both in redis and db, will set its key DISABLE
	// and return 'fmt.Errorf("not found record in db nor redis")'
	if e == nil && string(buf) == "DISABLE" {
		return nil, fmt.Errorf("not found record in db nor redis")
	}

	// Not found in redis
	if e == redis.ErrNil {
		return nil, e
	}

	// Server error, should be logged by caller
	if e != nil && e != redis.ErrNil {
		return nil, errorx.Wrap(e)
	}

	return buf, nil
}

// no-decode Tail


// Auto-Generate Tail