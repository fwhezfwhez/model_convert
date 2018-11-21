package cmd

import (
	"cdd-platform-srv/utils/model_convert"
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/fwhezfwhez/errorx"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

var m sync.Mutex
var dataSource string
var tables []string
var exportTo string

func Init() {
	var tableNames string
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("conf")
	ReadConfig(v)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		ReadConfig(v)
	})

	// 命令行指定表,数据源
	flag.StringVar(&tableNames, "tables", "", "go run main.go -tables 'animal,bird,dog,cat'")
	flag.StringVar(&dataSource, "dataSource", "", "go run main.go -dataSource `host=localhost port=5432 user=postgres dbname=test sslmode=disable password=123`")
	flag.StringVar(&exportTo, "exportTo","", "go run main.go -exportTo '/tmp/model.go'")

	flag.Parse()
	tmp := ""
	if tableNames == "" {
		tmp = v.GetString("tableNames")
		if tmp == "" {
			panic("tables not specific neither by '-tables x,y,z' nor conf.yaml")
		} else {
			tableNames = tmp
		}
	}
	tables = strings.Split(tableNames, ",")


	if dataSource == ""{
		tmp = v.GetString("dataSource")
		if tmp == "" {
			panic("dataSource not specific neither by '-dataSource x,y,z' nor conf.yaml")
		} else {
			dataSource = tmp
		}
	}

	if exportTo == ""{
		tmp = v.GetString("exportTo")
		if tmp == "" {
			panic("exportTo not specific neither by '-exportTo <path>' nor conf.yaml")
		} else {
			exportTo = tmp
		}
	}
}
func main() {
	Init()
	var result string
	for _,table :=range tables {
		result += model_convert.TableToStructWithTag(dataSource, table)
	}

}

func ReadConfig(v *viper.Viper) error {
	m.Lock()
	defer m.Unlock()
	err := v.ReadInConfig()
	if err != nil {
		return errorx.NewFromString("Error on parsing config file!")
	}
	return nil
}
