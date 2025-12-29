package cfg

import (
	"dos/logger"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

const (
	defaultPort      = "8080"
	defaultDBName    = "appdb"
	defaultDbUrl     = "localhost"
	defaultDbPort    = "5432"
	defaultDbUser    = "admin"
	defaultDbPwd     = "test"
	defaultFECorsUrl = "http://localhost:5173"
)

type Config struct {
	LogLevel  string
	Dsn       string
	MaskedDsn string
	AppPort   string
	FECorsUrl string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	logLvl := os.Getenv("LOG_LEVEL")
	logger.InitLogger(logLvl)

	dsnArray := getDsnMaskedDsn()
	dsn := dsnArray[0]
	maskedDsn := dsnArray[1]
	logger.L.Info("passed from config", "dsn", maskedDsn)

	appPort := getWithDefault("APP_PORT", defaultPort)
	feUrl := getWithDefault("FE_CORS_URL", defaultFECorsUrl)
	return &Config{Dsn: dsn, MaskedDsn: maskedDsn, LogLevel: logLvl, AppPort: appPort, FECorsUrl: feUrl}
}

func getDsnMaskedDsn() []string {
	url := getWithDefault("PG_DB_URL", defaultDbUrl)
	port := getWithDefault("PG_DB_PORT", defaultDbPort)
	user := getWithDefault("PG_DB_USERNAME", defaultDbUser)
	pwd := getWithDefault("PG_DB_PASSWORD", defaultDbPwd)
	db := getWithDefault("PG_DB_NAME", defaultDBName)
	maskedDsn := fmt.Sprintf("postgres://******:******@%s:%s/%s?sslmode=disable", url, port, db)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pwd, url, port, db)

	return []string{dsn, maskedDsn}
}

func getWithDefault(env, defaultVal string) string {
	value := os.Getenv(env)
	if value == "" {
		logger.L.Warn("environment variable not found", "key", env, "default", defaultVal)
		value = defaultVal
	}
	logger.L.Info("environment variable will be used", "key", env, "value", value)
	return value
}
