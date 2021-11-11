package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"red_envelope/config"
	"red_envelope/utils"

	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type App struct {
	RDB              *redis.Client
	EnvelopeProducer *Producer
	MaxCount         int   // 每个uid最多抢到的红包数
	MaxAmount        int64 // 设置的红包总金额
	MaxSize          int64 // 设置的红包总数量
	SnatchedPr       int   // 抢到红包概率
	RemainingAmount  int64 // 可发红包总金额
	RemainingSize    int64 // 可发红包总数
	UserCount        *cache.Cache
	UserWallet       *cache.Cache
	KafkaProducer    *KafkaProducer
	UserMutex        sync.Map
}

var (
	onceApp *App
	once    sync.Once
	ctx     = context.Background()
)

func GetApp() *App {
	once.Do(func() {
		onceApp = &App{
			UserCount:  cache.New(5*time.Minute, 10*time.Minute),
			UserWallet: cache.New(5*time.Minute, 10*time.Minute),
		}
	})
	return onceApp
}

func (app *App) Run() {
	// 数据库连接
	app.OpenRedis()
	app.OpenKafkaProducer()
	go app.KafkaProducer.HandleSendErr()
	// 参数配置加载
	app.LoadConfig()

	// 开始生产红包
	app.EnvelopeProducer = NewProducer(app.RemainingAmount, app.RemainingSize)
	go app.EnvelopeProducer.Do()
	app.EnvelopeProducer.MsgChan <- 1
}

func (app *App) OpenKafkaProducer() {
	kafkaBrokers := utils.GetEnv("KAFKA_ADDRS", config.DefaultKafkaBrokers)
	brokers := utils.GetArgs(kafkaBrokers)
	topic := utils.GetEnv("KAFKA_TOPIC", config.DefaultKafkaTopic)
	kafkaProducer := GetKafkaProducer(topic, brokers)
	app.KafkaProducer = &kafkaProducer
	logrus.Infoln("success connect to Kafka")
}

func (app *App) OpenRedis() {

	host := utils.GetEnv("REDIS_SERVICE_HOST", config.DefaultHost)
	port := utils.GetEnv("REDIS_SERVICE_PORT", config.DefaultRedisPort)
	password := utils.GetEnv("REDIS_PASSWORD", config.DefaultRedisPasswd)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	app.RDB = rdb
	logrus.Infoln("success connect to redis")
}

func (app *App) LoadConfig() {
	var err error
	var curAmount, curSize int64
	// max_amount 先从redis中取，没有则从env中初始化
	var maxAmount, maxSize string
	val, err := app.RDB.Get(ctx, "max_amount").Result()
	if err != nil {
		logrus.Info("load max_amount from env...")
		maxAmount = utils.GetEnv("MAX_AMOUNT", config.DefaultMaxAmount)
		if err := app.RDB.Set(ctx, "max_amount", maxAmount, 0).Err(); err != nil {
			logrus.Fatalln("write max_amount to redis failed...")
		}
	} else {
		logrus.Info("load max_amount from redis...")
		maxAmount = val
	}
	if app.MaxAmount, err = strconv.ParseInt(maxAmount, 10, 64); err != nil {
		logrus.Fatalln("load max_amount failed...", err)
	}

	if curAmount, err = app.GetCurAmount(); err != nil {
		if err := app.RDB.Set(ctx, "cur_amount", 0, 0).Err(); err != nil {
			logrus.Fatalln("load max_amount failed...", err)
		}
	}
	app.RemainingAmount = app.MaxAmount - curAmount

	maxCount := utils.GetEnv("MAX_COUNT", config.DefaultMaxCount)
	if app.MaxCount, err = strconv.Atoi(maxCount); err != nil {
		logrus.Fatalln("load max_count failed...", err)
	}
	// max_size 先从redis中取，没有则从env中初始化
	val, err = app.RDB.Get(ctx, "max_size").Result()
	if err != nil {
		logrus.Info("load max_size from env...")
		maxSize = utils.GetEnv("MAX_SIZE", config.DefaultMaxSize)
		if err := app.RDB.Set(ctx, "max_size", maxSize, 0).Err(); err != nil {
			logrus.Fatalln("write max_size to redis failed...")
		}
	} else {
		logrus.Info("load max_size from redis...")
		maxSize = val
	}
	if app.MaxSize, err = strconv.ParseInt(maxSize, 10, 64); err != nil {
		logrus.Fatalln("load max_size failed...", err)
	}

	if curSize, err = app.GetCurSize(); err != nil {
		if err := app.RDB.Set(ctx, "cur_size", 0, 0).Err(); err != nil {
			logrus.Fatalln("load max_size failed...", err)
		}
	}
	app.RemainingSize = app.MaxSize - curSize

	var ok bool
	snatchedPr := utils.GetEnv("SNATCHED_PR", config.DefaultSnatchedPr)
	if app.SnatchedPr, ok = CheckSnatchedPr(snatchedPr); !ok {
		logrus.Fatalln("load snatched_pr failed...", err)
	}

	logrus.Infoln("success load config")
}

func (app *App) GetCurAmount() (curAmount int64, err error) {
	var val string
	if val, err = app.RDB.Get(ctx, "cur_amount").Result(); err != nil {
		return 0, err
	}
	if curAmount, err = strconv.ParseInt(val, 10, 64); err != nil {
		return 0, err
	}
	return curAmount, err
}

func (app *App) GetCurSize() (curSize int64, err error) {
	var val string
	if val, err = app.RDB.Get(ctx, "cur_size").Result(); err != nil {
		return 0, err
	}
	if curSize, err = strconv.ParseInt(val, 10, 64); err != nil {
		return 0, err
	}
	return curSize, err
}

func (app *App) AddAmount(val int64) {
	app.MaxAmount += val
	app.RemainingAmount += val
	app.EnvelopeProducer.Mutex.Lock()
	app.EnvelopeProducer.Amount += val
	app.EnvelopeProducer.Mutex.Unlock()
}

func (app *App) RollbackAddAmount(val int64) {
	app.MaxAmount -= val
	app.RemainingAmount -= val
	app.EnvelopeProducer.Mutex.Lock()
	app.EnvelopeProducer.Amount -= val
	app.EnvelopeProducer.Mutex.Unlock()
}

func (app *App) AddSize(val int64) {
	app.MaxSize += val
	app.RemainingSize += val
	app.EnvelopeProducer.Mutex.Lock()
	app.EnvelopeProducer.Size += val
	app.EnvelopeProducer.Mutex.Unlock()
}

func (app *App) RollbackAddSize(val int64) {
	app.MaxSize -= val
	app.RemainingSize -= val
	app.EnvelopeProducer.Mutex.Lock()
	app.EnvelopeProducer.Size -= val
	app.EnvelopeProducer.Mutex.Unlock()
}

func CheckSnatchedPr(snatchedPr string) (value int, ok bool) {
	var err error
	if value, err = strconv.Atoi(snatchedPr); err != nil {
		return
	}
	if value >= 0 && value <= 100 {
		ok = true
	}
	return
}
