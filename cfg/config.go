package cfg

import (
	"dos/logger"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Config struct {
	LogLevel  string
	Dsn       string
	MaskedDsn string
	AppPort   string
}

func getDsnMaskedDsn() []string {

	url := os.Getenv("PG_DB_URL")
	port := os.Getenv("PG_DB_PORT")
	user := os.Getenv("PG_DB_USERNAME")
	pwd := os.Getenv("PG_DB_PASSWORD")
	db := os.Getenv("PG_DB_NAME")
	maskedDsn := fmt.Sprintf("postgres://******:******@%s:%s/%s?sslmode=disable", url, port, db)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pwd, url, port, db)

	return []string{dsn, maskedDsn}
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

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		slog.Error("APP_PORT not found in configuration. Default 8080 is used")
		appPort = "8080"
	}
	return &Config{Dsn: dsn, MaskedDsn: maskedDsn, LogLevel: logLvl, AppPort: appPort}
}
