package initialize

import (
	"context"
	"fmt"
	"log"
	"red_packet/router"
	"red_packet/utils"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	DB  *gorm.DB
	RDB *redis.Client
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
	app.openDB()
	app.openRedis()

	r := router.InitRouter()
	log.Fatal(r.Listen(":8080"))
}

func (app *App) openDB() {
	host := utils.GetEnv("MYSQL_SERVICE_HOST", "localhost")
	port := utils.GetEnv("MYSQL_SERVICE_PORT", "3306")
	password := utils.GetEnv("MYSQL_ROOT_PASSWORD", "123456")

	dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/test?charset=utf8mb4&parseTime=True&loc=Local", password, host, port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	app.DB = db
	log.Println("success connect to mysql")
}

func (app *App) openRedis() {
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
	log.Println("success connect to redis")
}
