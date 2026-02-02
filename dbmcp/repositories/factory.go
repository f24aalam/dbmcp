package repositories

func GetRepository(dbType string) TableRepository {
	switch dbType {
	case "sqlite":
		return NewSQLiteRepository()
	case "mysql":
		return NewMySQLRepository()
	default:
		return NewMySQLRepository()
	}
}
