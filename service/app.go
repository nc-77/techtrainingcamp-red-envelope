package service

import (
	"context"
	"fmt"
	"red_packet/config"
	"strconv"
	"sync"

	"red_packet/utils"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	DB               *gorm.DB
	RDB              *redis.Client
	EnvelopeProducer *Producer
	MaxCount         int   // 每个uid最多抢到的红包数
	MaxAmount        int64 // 红包总金额
	MaxSize          int64 // 红包总数量
	UserCount        sync.Map
}

var (
	onceApp *App
	once    sync.Once
	ctx     = context.Background()
)

func GetApp() *App {
	once.Do(func() {
		onceApp = &App{}
	})
	return onceApp
}

func (app *App) Run() {
	// 数据库连接
	app.OpenDB()
	app.OpenRedis()

	// 参数配置加载
	app.LoadConfig()

	// 开始生产红包
	app.EnvelopeProducer = NewProducer(app.MaxAmount, app.MaxSize)
	go app.EnvelopeProducer.Do()
}

func (app *App) OpenDB() {
	host := utils.GetEnv("MYSQL_SERVICE_HOST", config.DefaultHost)
	port := utils.GetEnv("MYSQL_SERVICE_PORT", config.DefaultMySQLPort)
	password := utils.GetEnv("MYSQL_ROOT_PASSWORD", config.DefaultMySQLPasswd)
	dbName := utils.GetEnv("MYSQL_DB", config.DefaultMySQLDB)

	dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", password, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	app.DB = db
	logrus.Infoln("success connect to mysql")
}

func (app *App) OpenRedis() {
	ctx := context.Background()
	host := utils.GetEnv("REDIS_SERVICE_HOST", config.DefaultHost)
	port := utils.GetEnv("REDIS_SERVICE_PORT", config.DefaultRedisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: config.DefaultRedisPasswd, // no password set
		DB:       0,                         // use default DB
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	app.RDB = rdb
	logrus.Infoln("success connect to redis")
}

func (app *App) LoadConfig() {
	var err error
	amount := utils.GetEnv("AMOUNT", config.DefaultMaxAmount)
	if app.MaxAmount, err = strconv.ParseInt(amount, 10, 64); err != nil {
		logrus.Fatalln("load amount failed...")
	}
	maxCount := utils.GetEnv("MAX_COUNT", config.DefaultMaxCount)
	if app.MaxCount, err = strconv.Atoi(maxCount); err != nil {
		logrus.Fatalln("load max_count failed...")
	}
	maxSize := utils.GetEnv("MAX_SIZE", config.DefaultMaxSize)
	if app.MaxSize, err = strconv.ParseInt(maxSize, 10, 64); err != nil {
		logrus.Fatalln("load max_size failed...")
	}
	logrus.Infoln("success load config")
}