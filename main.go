package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

func getEnv(env string, defaultVal string) (key string) {
	if key = os.Getenv(env); key == "" {
		key = defaultVal
	}
	return
}

var (
	ctx = context.Background()
)

func main() {
	// test mysql
	host := getEnv("MYSQL_SERVICE_HOST", "localhost")
	port := getEnv("MYSQL_SERVICE_PORT", "3306")
	password := getEnv("MYSQL_ROOT_PASSWORD", "123456")
	dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/test", password, host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	// test redis
	host = getEnv("REDIS_SERVICE_HOST", "localhost")
	port = getEnv("REDIS_SERVICE_PORT", "6379")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if _, err = rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	// http server
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!")
	})

	panic(app.Listen(":8080"))
}
