package src

import (
	"fmt"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/go-mysql-org/go-mysql/canal"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	once      sync.Once
	appConfig unsafe.Pointer
	DB        *gorm.DB
	RedisCli  *redis.Client
	CanalCfg  *canal.Config
)

type RedisCfg struct {
	Uri             string        `yaml:"uri"`
	ConnectTimeout  time.Duration `yaml:"connect_timeout" mapstructure:"connect_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
	PoolMaxActive   int           `yaml:"pool_max_active" mapstructure:"pool_max_active"`
	PoolIdleTimeout time.Duration `yaml:"pool_idle_timeout" mapstructure:"pool_idle_timeout"`
	MaxRetries      int           `yaml:"max_retries" mapstructure:"max_retries"`
	DB              int           `yaml:"db"`
}

func (opt RedisCfg) MakeRedisClient() interface{} {
	return RedisCli
}

type MySQLCfg struct {
	DSN         string        `yaml:"dsn"`
	MaxLifetime time.Duration `yaml:"max_lifetime" mapstructure:"max_lifetime"`
	MaxOpenConn int           `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConn int           `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	Autocommit  bool          `yaml:"autocommit"`
}

type CanalConfig struct {
	Addr     string   `yaml:"addr"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Flavor   string   `yaml:"flavor"`
	DB       string   `yaml:"db"`
	Table    []string `yaml:"table"`
	ServerId uint32   `yaml:"server_id" mapstructure:"server_id"`
}

type AppCfg struct {
	Redis *RedisCfg
	MySQL *MySQLCfg
	Canal *CanalConfig
}

func ResourceInit() {
	var must = func(err error) {
		if err != nil {
			panic(err)
		}
	}
	once.Do(func() {
		var err error
		// initialize mysql connection
		cfg := getAppCfg()
		DB, err = openMySQL(cfg.MySQL)
		must(err)
		// initialize redis connection
		RedisCli, err = openRedis(cfg.Redis)
		must(err)
		CanalCfg = createCanalCfg(cfg.Canal)
	})
}

func readAppCfg() {
	// app.yaml locates in the src directory
	fileLocation()
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Sprintf("app.yaml not found:%s", err.Error()))
		} else {
			panic(fmt.Sprintf("app.yaml read error:%s", err.Error()))
		}
	}
	var cfg AppCfg
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Sprintf("unmarshal config fail:%s", err.Error()))
	}
	cfgStr, _ := jsoniter.MarshalToString(cfg)
	logrus.Infof("cfg:%s", cfgStr)
	atomic.StorePointer(&appConfig, unsafe.Pointer(&cfg))
}

// find directory where app.yaml locates
func fileLocation() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	parentDir := pwd
	var appYamlDir string
	for {
		if judgeFileExist(path.Join(parentDir, "app.yaml")) {
			appYamlDir = parentDir
			break
		}
		temp := path.Dir(parentDir)
		if parentDir == temp {
			break
		}
		parentDir = temp
	}
	if appYamlDir != "" {
		_ = os.Chdir(appYamlDir)
	}
	return appYamlDir
}

func judgeFileExist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if os.IsExist(err) {
		return true
	}
	return false
}

func getAppCfg() *AppCfg {
	if appConfig == nil {
		readAppCfg()
	}
	p := atomic.LoadPointer(&appConfig)
	return (*AppCfg)(p)
}

func openMySQL(dbConfig *MySQLCfg) (*gorm.DB, error) {
	dsn := attachDBDsn(dbConfig.DSN, dbConfig.Autocommit, false)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(dbConfig.MaxLifetime)
	err = sqlDB.Ping()
	return db, err
}

func attachDBDsn(dsn string, autocommit, isSlave bool) string {
	dsn += "?interpolateParams=true&parseTime=true&loc=Local&autocommit="
	if isSlave {
		dsn += "1"
	} else {
		if autocommit {
			dsn += "1"
		} else {
			dsn += "0"
		}
	}
	return dsn
}

func openRedis(redisCfg *RedisCfg) (*redis.Client, error) {
	option := &redis.Options{
		Addr:         redisCfg.Uri,
		DialTimeout:  redisCfg.ConnectTimeout,
		ReadTimeout:  redisCfg.ReadTimeout,
		WriteTimeout: redisCfg.WriteTimeout,
		PoolSize:     redisCfg.PoolMaxActive,
		MaxRetries:   redisCfg.MaxRetries,
		DB:           redisCfg.DB,
	}
	client := redis.NewClient(option)
	return client, nil
}

func createCanalCfg(canalCfg *CanalConfig) *canal.Config {
	return &canal.Config{
		Addr:     canalCfg.Addr,
		User:     canalCfg.User,
		Password: canalCfg.Password,
		Flavor:   canalCfg.Flavor,
		Dump: canal.DumpConfig{
			Tables:  canalCfg.Table,
			TableDB: canalCfg.DB,
		},
		ServerID: canalCfg.ServerId,
	}
}
