package config

// AppConfig is a struct that holds the configuration of the app
type AppConfig struct {
	Env *EnvVariables
}

// EnvVariables is a struct that holds the environment variables
type EnvVariables struct {
	DSN string
}
