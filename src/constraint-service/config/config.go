package config

import (
	"os"
	"strconv"
)

type Config struct {
	HeartbeatPort      int
	GrpcPort           int
	UseReflection      bool
	DbConnectionString string
	UseStartupMigrate  bool
	MigrationDirectory string
}

func LoadConfig() *Config {
	grpc_port, err := strconv.Atoi(os.Getenv("CS_GRPC_PORT"))
	if err != nil {
		grpc_port = 8001
	}

	heartbeat_port, err := strconv.Atoi(os.Getenv("CS_HEART_PORT"))
	if err != nil {
		heartbeat_port = 8002
	}

	return &Config{
		HeartbeatPort:      heartbeat_port,
		GrpcPort:           grpc_port,
		DbConnectionString: os.Getenv("CS_DB_STRING"),
		UseReflection:      os.Getenv("CS_REFLECTION") == "true",
		UseStartupMigrate:  os.Getenv("CS_STARTUP_MIGRATE") == "true",
		MigrationDirectory: os.Getenv("CS_MIGRATION_DIR"),
	}
}
