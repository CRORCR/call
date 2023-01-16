package model

type Config struct {
	App   AppConfig   `yaml:"app"`   // app相关配置
	Log   LogConfig   `yaml:"log"`   // log日志配置
	Redis RedisConfig `yaml:"redis"` // 缓存配置
	Mysql MysqlConfig `yaml:"mysql"` // mysql配置 测试先用着，链路调通之后，换成pgsql
}

type App struct {
	Env  string `yaml:"env"`
	Port uint64 `yaml:"port"`
}

// AppConfig app.yaml
type AppConfig struct {
	Env  string `yaml:"env"`
	Port string `yaml:"port"`
}

// LogConfig log.yaml
type LogConfig struct {
	level      string `yaml:"level"`
	MaxAge     string `yaml:"max_age"`
	DateFormat string `yaml:"date_format"`
}

// RedisConfig redis.yaml
type RedisConfig struct {
	Host           string `yaml:"host"`
	Port           int64  `yaml:"port"`
	Password       string `yaml:"password"`
	MaxActive      int64  `yaml:"max_active"`
	IdleTimeout    int64  `yaml:"idle_timeout"`
	ConnectTimeout int64  `yaml:"connect_timeout"`
	ReadTimeout    int64  `yaml:"read_timeout"`
	WriteTimeout   int64  `yaml:"write_timeout"`
	DB             int64  `yaml:"db"`
	Retry          int64  `yaml:"retry"`
}

// MysqlConfig mysql.yaml 拆分出默认的配置
type MysqlConfig struct {
	DefaultMaster DefaultDbConfig `yaml:"default_master"`
	DefaultSlave  DefaultDbConfig `yaml:"default_slave"`
	CallDb        DbConfig        `yaml:"call_db"`
}

type DefaultDbConfig struct {
	Driver       string `yaml:"driver"`        // 连接驱动
	Dsn          string `yaml:"dsn"`           // dsn，如果设置了dsn, 以下的所有设置都不生效
	Host         string `yaml:"host"`          // ip地址
	Port         int64  `yaml:"port"`          // 端口
	Username     string `yaml:"username"`      // 用户名
	Password     string `yaml:"password"`      // 密码
	Charset      string `yaml:"charset"`       // 字符集
	Collation    string `yaml:"collation"`     // 字符序
	Timeout      int64  `yaml:"timeout"`       // 连接超时
	ReadTimeout  int64  `yaml:"read_timeout"`  // 读超时 单位s
	WriteTimeout int64  `yaml:"write_timeout"` // 写超时 单位s
	ParseTime    bool   `yaml:"parse_time"`    // 是否解析时间
	Protocol     string `yaml:"protocol"`      // 传输协议
	Loc          string `yaml:"loc"`           // 时区
}

// DbConfig mysql分主从两份配置，一期先用单库
type DbConfig struct {
	Master DatabaseConfig `yaml:"master"`
	Slave  DatabaseConfig `yaml:"slave"`
}

type DatabaseConfig struct {
	Database              string `yaml:"database"`                 // 数据库名
	MaxIdleConnections    int64  `yaml:"max_idle_connections"`     // 最大空闲连接数，默认 100
	MaxOpenConnections    int64  `yaml:"max_open_connections"`     // 最大打开的连接数，默认 100
	MaxConnectionLifeTime int64  `yaml:"max_connection_life_time"` // 空闲连接最大存活时间，默认 10s
}
