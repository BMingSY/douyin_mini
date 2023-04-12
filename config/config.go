package config

var Config Service

type Service struct {
	MySQL *MySQLConfig
	Redis *RedisConfig
	MinIO *MinIOConfig
	Mongo *MongoConfig
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
	Timeout  string `mapstructure:"timeout"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type MinIOConfig struct {
	Url             string `mapstructure:"url"`
	APIPort         string `mapstructure:"api_port"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

type MongoConfig struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Database   string `mapstructure:"database"`
	Collection string `mapstructure:"collection"`
}
