package postgres

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/adohong4/driving-license/config"
	"github.com/jmoiron/sqlx"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// Return new Postgresql db instance
func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	dataSourceName := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dataSourceName != "" {
		if strings.Contains(dataSourceName, "<") || strings.Contains(dataSourceName, ">") {
			return nil, errors.New("DATABASE_URL contains placeholder characters '<' or '>', please replace with real Render credentials")
		}
		if !strings.Contains(dataSourceName, "sslmode=") {
			if strings.Contains(dataSourceName, "?") {
				dataSourceName += "&sslmode=require"
			} else {
				dataSourceName += "?sslmode=require"
			}
		}
		db, err := sqlx.Connect(c.Postgres.PgDriver, dataSourceName)
		if err != nil {
			return nil, err
		}

		db.SetMaxOpenConns(maxOpenConns)
		db.SetConnMaxLifetime(connMaxLifetime * time.Second)
		db.SetMaxIdleConns(maxIdleConns)
		db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
		if err = db.Ping(); err != nil {
			return nil, err
		}
		return db, nil
	}

	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	password := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	dbName := strings.TrimSpace(os.Getenv("DB_NAME"))
	sslMode := strings.TrimSpace(os.Getenv("DB_SSLMODE"))
	if sslMode == "" {
		sslMode = strings.TrimSpace(os.Getenv("PGSSLMODE"))
	}

	dataSourceName = ""
	if host != "" && port != "" && user != "" && password != "" && dbName != "" {
		if strings.Contains(password, "<") || strings.Contains(password, ">") {
			return nil, errors.New("DB_PASSWORD contains placeholder characters '<' or '>', please replace with real Render credentials")
		}
		if sslMode == "" {
			sslMode = "disable"
			if c.Postgres.PostgresqlSSLMode {
				sslMode = "require"
			}
		}
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			host,
			port,
			user,
			dbName,
			sslMode,
			password,
		)
	}

	if dataSourceName == "" {
		sslMode := "disable"
		if c.Postgres.PostgresqlSSLMode {
			sslMode = "require"
		}
		if envSSLMode := os.Getenv("PGSSLMODE"); envSSLMode != "" {
			sslMode = envSSLMode
		}
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			c.Postgres.PostgresqlHost,
			c.Postgres.PostgresqlPort,
			c.Postgres.PostgresqlUser,
			c.Postgres.PostgresqlDbname,
			sslMode,
			c.Postgres.PostgresqlPassword,
		)
	} else if !strings.Contains(dataSourceName, "sslmode=") {
		if strings.Contains(dataSourceName, "?") {
			dataSourceName += "&sslmode=require"
		} else {
			dataSourceName += "?sslmode=require"
		}
	}

	db, err := sqlx.Connect(c.Postgres.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
