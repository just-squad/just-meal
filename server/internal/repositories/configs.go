package repositories

import "fmt"

type DbType string

const (
	Postgres DbType = "postgres"
	Mongo    DbType = "mongo"
)

// Конфигурация БД
type DBConfig struct {
	Type     DbType      // Тип базы данных
	Postgres PgConfig    // Конфигурация подключения к бд Postgres
	Mongo    MongoConfig // Конфигурация подключения к бд MongoDb
}

type PgConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (c *PgConfig) GetConnectionString() string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
	return psqlInfo
}

type MongoConfig struct {
	URI      string // Адрес подключения к MongoDb
	Database string // Название базы данных
}
