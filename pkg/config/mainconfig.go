package config

import (
	"github.com/go-xorm/xorm"
)

// GeneralConfig has miscellaneous configuration options
type GeneralConfig struct {
	InstanceID string `toml:"instanceID"`
	LogDir     string `toml:"logdir"`
	HomeDir    string `toml:"homedir"`
	DataDir    string `toml:"datadir"`
	LogLevel   string `toml:"loglevel"`
}

//DatabaseCfg de configuration for the database
type DatabaseCfg struct {
	Type       string `toml:"type"`
	Host       string `toml:"host"`
	Name       string `toml:"name"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	SQLLogFile string `toml:"sqllogfile"`
	Debug      string `toml:"debug"`
	x          *xorm.Engine
	numChanges int64 `toml:"-"`
}

//SelfMonConfig configuration for self monitoring
type SelfMonConfig struct {
	Enabled           bool     `toml:"enabled"`
	Freq              int      `toml:"freq"`
	Prefix            string   `toml:"prefix"`
	InheritDeviceTags bool     `toml:"inheritdevicetags"`
	ExtraTags         []string `toml:"extra-tags"`
}

//HTTPConfig has webserver config options
type HTTPConfig struct {
	Port          int    `toml:"port"`
	AdminUser     string `toml:"adminuser"`
	AdminPassword string `toml:"adminpassword"`
	CookieID      string `toml:"cookieid"`
}

type InfluxCfg struct {
	ID        string `toml:"id"`
	Host      string `toml:"host"`
	Port      int    `toml:"port"`
	DB        string `toml:"db"`
	User      string `toml:"user"`
	Password  string `toml:"password"`
	Retention string `toml:"retention"`
	Precision string `toml:"precision"` //posible values [h,m,s,ms,u,ns] default seconds for the nature of data
	Timeout   int    `toml:"timeout"`
	UserAgent string `toml:"useragent"`
}

type Config struct {
	General  GeneralConfig
	Database DatabaseCfg
	Selfmon  SelfMonConfig
	HTTP     HTTPConfig
	Influxdb InfluxCfg
}

//var MainConfig Config
