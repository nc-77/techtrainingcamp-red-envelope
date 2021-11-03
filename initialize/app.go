package initialize

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"red_packet/utils"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	DB       *gorm.DB
	RDB      *redis.Client
	MaxCount int    //每个uid最多抢到的红宝书
	Amount   uint64 // 红包总金额
}

var (
	onceApp *App
	once    sync.Once
)

func NewApp() *App {
	once.Do(func() {
		onceApp = &App{}
	})
	return onceApp
}

func (app *App) Run() {
	app.OpenDB()
	app.OpenRedis()
	app.LoadConfig()
}

func (app *App) OpenDB() {
	host := utils.GetEnv("MYSQL_SERVICE_HOST", "localhost")
	port := utils.GetEnv("MYSQL_SERVICE_PORT", "3306")
	password := utils.GetEnv("MYSQL_ROOT_PASSWORD", "123456")

	dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/test?charset=utf8mb4&parseTime=True&loc=Local", password, host, port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	app.DB = db
	logrus.Infoln("success connect to mysql")
}

func (app *App) OpenRedis() {
	ctx := context.Background()
	host := utils.GetEnv("REDIS_SERVICE_HOST", "localhost")
	port := utils.GetEnv("REDIS_SERVICE_PORT", "6379")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	app.RDB = rdb
	logrus.Infoln("success connect to redis")
}

func (app *App) LoadConfig() {
	var err error
	amount := utils.GetEnv("AMOUNT", "10000")
	if app.Amount, err = strconv.ParseUint(amount, 10, 64); err != nil {
		log.Fatalln("load amount failed...")
	}
	maxCount := utils.GetEnv("MAX_COUNT", "10")
	if app.MaxCount, err = strconv.Atoi(maxCount); err != nil {
		log.Fatalln("load max_count failed...")
	}

	logrus.Infoln("success load config")
}
