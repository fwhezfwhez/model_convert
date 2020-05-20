package model

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
	"testing"
	"time"
)

// 获取某个用户信息
func TestMustGet(t *testing.T) {
	rp, db := initConfig()

	conn := rp.Get()
	defer conn.Close()

	var ui = UserInfo{
		UserId: 10086,
		GameId: 78,
	}

	engine := db.Model(&ui).Where("game_id=? and user_id=?", ui.GameId, ui.UserId)

	if e := ui.MustGet(conn, engine); e != nil {
		if e != redis.ErrNil && e.Error() != "not found record in db nor redis" {
			panic(e)
		}
		fmt.Println("not such record")
		return
	}

	fmt.Println(ui)
}

// 返回某个用户信息，无decode开销版本
// 返回的[]byte值可能为空，非空时表示在redis中找到，为空时表示该记录已经被数据库找到，并被赋予了ui本身
func TestMustGetNoDecode(t *testing.T) {
	rp, db := initConfig()

	conn := rp.Get()
	defer conn.Close()

	var ui = UserInfo{
		UserId: 10086,
		GameId: 78,
	}

	engine := db.Model(&ui).Where("game_id=? and user_id=?", ui.GameId, ui.UserId)

	buf, e := ui.MustGetNoDecode(conn, engine)
	if e != nil {
		if e != redis.ErrNil && e.Error() != "not found record in db nor redis" {
			panic(e)
		}
		fmt.Println("not such record")
		return
	}

	if len(buf) > 0 {
		fmt.Println(string(buf))
	} else {
		fmt.Println(ui)
	}
}

// 返回某个用户信息，无decode开销版本
// 返回无json-decode开销的单条数据
func TestGetFromRedisNoDecode(t *testing.T) {
	rp, _ := initConfig()

	conn := rp.Get()
	defer conn.Close()

	var ui = UserInfo{
		UserId: 10086,
		GameId: 78,
	}

	buf, e := ui.GetFromRedisNoDecode(conn)
	if e != nil {
		if e != redis.ErrNil && e.Error() != "not found record in db nor redis" {
			panic(e)
		}
		fmt.Println("not such record")
		return
	}

	fmt.Println(string(buf))
}

// 获取某个game_id的所有用户
func TestArrayMustGet(t *testing.T) {
	rp, db := initConfig()

	conn := rp.Get()
	defer conn.Close()

	var ui = UserInfo{
		GameId: 78,
	}

	engine := db.Model(&ui).Where("game_id=?", ui.GameId)

	arr, e := ui.ArrayMustGet(conn, engine)
	if e != nil {
		if e != redis.ErrNil && e.Error() != "not found record in db nor redis" {
			panic(e)
		}
		fmt.Println("not such record")
		return
	}
	fmt.Println(len(arr))

	fmt.Println(arr)
}

// 获取某个game_id的所有用户，无decode开销版本
// 返回的[]byte值可能为空，非空时表示在redis中找到，为空时表示该记录已经被数据库找到，并被赋予了o本身
func TestUserInfo_ArrayMustGetNoDecode(t *testing.T) {
	rp, db := initConfig()

	conn := rp.Get()
	defer conn.Close()

	var ui = UserInfo{
		GameId: 78,
	}

	engine := db.Model(&ui).Where("game_id=?", ui.GameId)

	arr, buf, e := ui.ArrayMustGetNoDecode(conn, engine)
	if e != nil {
		if e != redis.ErrNil && e.Error() != "not found record in db nor redis" {
			panic(e)
		}
		fmt.Println("not such record")
		return
	}
	if len(buf) == 0 {
		fmt.Println(len(arr))
	} else {
		fmt.Println(len(arr))
		fmt.Println(string(buf))
	}
}

// 获取某个game_id的所有用户，无decode开销版本
func TestUserInfo_ArrayGetFromRedisNoDecode(t *testing.T) {
	rp, _ := initConfig()
	conn := rp.Get()
	defer conn.Close()

	ui := UserInfo{
		GameId: 78,
	}

	buf, e := ui.ArrayGetFromRedisNoDecode(conn)
	if e != nil {
		if e != redis.ErrNil && e.Error() != "not found record in db nor redis" {
			panic(e)
		}
		fmt.Println("not such record")
		return
	}

	fmt.Println(string(buf))
}

func initConfig() (*redis.Pool, *gorm.DB) {
	// db
	dbConfig := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		"localhost",
		"postgres",
		"game",
		"disable",
		"123",
	)
	fmt.Println(dbConfig)
	db, err := gorm.Open("postgres",
		dbConfig,
	)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetConnMaxLifetime(10 * time.Second)
	db.DB().SetMaxIdleConns(1)
	if e := db.DB().Ping(); e != nil {
		panic(e)
	}

	// redis pool
	newPool := func(server, password string, db int) *redis.Pool {
		return &redis.Pool{
			MaxIdle:     200,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", server)
				if err != nil {
					fmt.Printf("occur error at newPool Dial: %v\n", err)
					return nil, err
				}
				_, e := c.Do("ping")
				if e != nil {
					if password != "" {
						if _, err := c.Do("AUTH", password); err != nil {
							c.Close()
							fmt.Printf("occur error at newPool Do Auth: %v\n", err)
							return nil, err
						}
						if _, e := c.Do("ping"); e != nil {
							return nil, fmt.Errorf("ping twice err: %v", e)
						}
					}
				}

				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					fmt.Printf("occur error at newPool Do SELECT: %v\n", err)
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}

	redisPool := newPool("localhost:6379", "", 0)
	return redisPool, db
}
