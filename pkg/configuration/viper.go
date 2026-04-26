package configuration

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
	Database  DatabaseConfig  `mapstructure:"database"`
	Messaging MessagingConfig `mapstructure:"messaging"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	App       AppConfig       `mapstructure:"app"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	PDF       PDFConfig       `mapstructure:"pdf"`
	Excel     ExcelConfig     `mapstructure:"excel"`
	Email     EmailConfig     `mapstructure:"email"`
}

// DatabaseConfig holds database configurations
type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

// PoolConfig holds database connection pool configuration
type PoolConfig struct {
	Idle     int `mapstructure:"idle"`
	Max      int `mapstructure:"max"`
	Lifetime int `mapstructure:"lifetime"`
}

// PostgresConfig holds PostgreSQL configuration
type PostgresConfig struct {
	Host      string     `mapstructure:"host"`
	Port      string     `mapstructure:"port"`
	Username  string     `mapstructure:"username"`
	Password  string     `mapstructure:"password"`
	DBName    string     `mapstructure:"dbname"`
	IsMigrate bool       `mapstructure:"is_migrate"`
	Pool      PoolConfig `mapstructure:"pool"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	URL      string     `mapstructure:"url"`
	Password string     `mapstructure:"password"`
	DB       int        `mapstructure:"db"`
	Pool     PoolConfig `mapstructure:"pool"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

// AppConfig holds application configuration
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Port    string `mapstructure:"port"`
	Prefork bool   `mapstructure:"prefork"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	UseLogrus            bool   `mapstructure:"use_logrus"`
	LogToConsole         bool   `mapstructure:"log_to_console"`
	LogToFile            bool   `mapstructure:"log_to_file"`
	LogSeparateByLevel   bool   `mapstructure:"log_separate_by_level"`
	LogCustomInfoFile    string `mapstructure:"log_custom_info_file"`
	LogCustomWarningFile string `mapstructure:"log_custom_warning_file"`
	LogCustomErrorFile   string `mapstructure:"log_custom_error_file"`
	LogDirectory         string `mapstructure:"log_directory"`
}

// MessagingConfig holds messaging configurations
type MessagingConfig struct {
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
}

// RabbitMQConfig holds RabbitMQ configuration
type RabbitMQConfig struct {
	URL      string `mapstructure:"url"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// KafkaProducerConfig holds Kafka producer configuration
type KafkaProducerConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers   []string            `mapstructure:"brokers"`
	Username  string              `mapstructure:"username"`
	Password  string              `mapstructure:"password"`
	Producer  KafkaProducerConfig `mapstructure:"producer"`
	Bootstrap BootstrapConfig     `mapstructure:"bootstrap"`
	Group     GroupConfig         `mapstructure:"group"`
	Auto      AutoConfig          `mapstructure:"auto"`
}

// BootstrapConfig holds Kafka bootstrap configuration
type BootstrapConfig struct {
	Servers string `mapstructure:"servers"`
}

// GroupConfig holds Kafka group configuration
type GroupConfig struct {
	ID string `mapstructure:"id"`
}

// AutoConfig holds Kafka auto configuration
type AutoConfig struct {
	Offset OffsetConfig `mapstructure:"offset"`
}

// OffsetConfig holds Kafka offset configuration
type OffsetConfig struct {
	Reset string `mapstructure:"reset"`
}

// PDFConfig holds PDF generation configuration
type PDFConfig struct {
	TemplateDir string `mapstructure:"template_dir"`
	BinaryPath  string `mapstructure:"binary_path"`
}

// ExcelConfig holds Excel configuration
type ExcelConfig struct {
	StorageDir string `mapstructure:"storage_dir"`
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost                string `mapstructure:"smtp_host"`
	SMTPPort                int    `mapstructure:"smtp_port"`
	SMTPUsername            string `mapstructure:"smtp_username"`
	SMTPPassword            string `mapstructure:"smtp_password"`
	FromEmail               string `mapstructure:"from_email"`
	FromName                string `mapstructure:"from_name"`
	TemplateDir             string `mapstructure:"template_dir"`
	BaseURL                 string `mapstructure:"base_url"`
	VerificationTokenExpiry int    `mapstructure:"verification_token_expiry"`
	ResetTokenExpiry        int    `mapstructure:"reset_token_expiry"`
}

var appConfig *Config

// LoadConfig loads configuration from config.json file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")         // Look for config file in the current directory
	viper.AddConfigPath("./")        // Look for config file in the current directory
	viper.AddConfigPath("../")       // Look for config file in the parent directory
	viper.AddConfigPath("../../")    // Look for config file in the grandparent directory
	viper.AddConfigPath("../../../") // Look for config file in the great-grandparent directory

	// Enable environment variable support
	viper.AutomaticEnv()

	// Allow viper to read Environment variables
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Set the global config
	appConfig = &config

	// log.Printf("Configuration loaded successfully from: %s", viper.ConfigFileUsed())
	return &config, nil
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	if appConfig == nil {
		config, err := LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
		return config
	}
	return appConfig
}

// Helper functions to access specific config values
func GetJWTSecret() string {
	return GetConfig().JWT.Secret
}

func GetPostgresConfig() PostgresConfig {
	return GetConfig().Database.Postgres
}

func GetRedisConfig() RedisConfig {
	return GetConfig().Database.Redis
}

func GetAppConfig() AppConfig {
	return GetConfig().App
}

func GetAppPort() string {
	return GetConfig().App.Port
}

func GetAppName() string {
	return GetConfig().App.Name
}

func GetAppPrefork() bool {
	return GetConfig().App.Prefork
}

func GetLoggingConfig() LoggingConfig {
	return GetConfig().Logging
}

func GetMessagingConfig() MessagingConfig {
	return GetConfig().Messaging
}

func GetRabbitMQConfig() RabbitMQConfig {
	return GetConfig().Messaging.RabbitMQ
}

func GetKafkaConfig() KafkaConfig {
	return GetConfig().Messaging.Kafka
}

func GetPDFConfig() PDFConfig {
	return GetConfig().PDF
}

func GetExcelConfig() ExcelConfig {
	return GetConfig().Excel
}

func GetEmailConfig() EmailConfig {
	return GetConfig().Email
}
