package config

// default用于本地测试环境
// const (
// 	DefaultHost = "localhost"

// 	DefaultMySQLPort   = "3306"
// 	DefaultMySQLPasswd = "123456"
// 	DefaultMySQLDB     = "test"

// 	DefaultRedisPort   = "6379"
// 	DefaultRedisPasswd = ""

// 	DefaultMaxCount  = "10"
// 	DefaultMaxAmount = "1000"
// 	DefaultMaxSize   = "1000"

// 	DefaultKafkaBrokers = "127.0.0.1:9092"
// 	DefaultKafkaTopic   = "test001"
// )

// lvpiche测试环境
const (
	DefaultHost = "localhost"

	DefaultMySQLPort     = "3306"
	DefaultMySQLUserName = "root"
	DefaultMySQLPasswd   = "123456"
	DefaultMySQLDB       = "test"

	DefaultRedisPort   = "6379"
	DefaultRedisPasswd = ""

	DefaultMaxCount   = "10"
	DefaultMaxAmount  = "1000"
	DefaultMaxSize    = "1000"
	DefaultSnatchedPr = "80" // 0-100

	DefaultKafkaBrokers = "127.0.0.1:9092"
	DefaultKafkaTopic   = "test001"
)
