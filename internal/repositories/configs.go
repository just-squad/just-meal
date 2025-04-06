package repositories

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
	URL string // Строка подключения к PostgresDb
}

type MongoConfig struct {
	URI      string // Адрес подключения к MongoDb
	Database string // Название базы данных
}
