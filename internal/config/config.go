package config

type AppConfig struct {
	Env *EnvVariables
}

type EnvVariables struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBPort     string
}
