package repositories

func GetRepository(dbType string) TableRepository {
	switch dbType {
	case "sqlite":
		return NewSQLiteRepository()
	case "mysql":
		return NewMySQLRepository()
	case "postgres":
		return NewPostgresRepository()
	default:
		return NewMySQLRepository()
	}
}
