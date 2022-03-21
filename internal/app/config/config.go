package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	GrpcPort               = ":9090"
	HttpPort               = ":8090"
	StdoutFileName         = "blog_system_out_log.log"
	StderrFileName         = "blog_system_out_log.log"
	TimestampFormat        = "02-01-2006 15:04:05"
	StdoutLoggerLevelValue = "debug"
	StderrLoggerLevelValue = "debug"
	FullTimestamp          = true
	WriteLoggerInfoInFile  = false

	DBHost     = "localhost"
	DBPort     = 54321
	DBUser     = "shem"
	DBPassword = "12345678"
	DBName     = "blog"
	DBSSLMode  = "disable"
)

type Config struct {
	Server     *server     `toml:"server"`
	LoggerConf *loggerConf `toml:"logger_conf"`
	DB         *db         `toml:"db"`
}

func NewConfig() *Config {
	return &Config{}
}

type server struct {
	GRPCPort string `toml:"port_grpc"`
	HTTPPort string `toml:"port_rest"`
	DocsPort string `toml:"port_docs"`
}

type loggerConf struct {
	StdoutFileName         string `toml:"stdout_file_name"`
	StderrFileName         string `toml:"stderr_file_name"`
	TimestampFormat        string `toml:"timestamp_format"`
	StdoutLoggerLevelValue string `toml:"stdout_logger_level"`
	StderrLoggerLevelValue string `toml:"stderr_logger_level"`
	FullTimestamp          bool   `toml:"full_timestamp"`
	WriteLoggerInfoInFile  bool   `toml:"write_logger_info_in_file"`
	StdoutFileWrite        *os.File
	StderrFileWrite        *os.File
	Formatter              *logrus.TextFormatter
	StdoutLoggerLevel      logrus.Level
	StderrLoggerLevel      logrus.Level
}

type db struct {
	Host     string `toml:"host"`
	Port     int64  `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
	SSLMode  string `toml:"ssl_mode"`
}
