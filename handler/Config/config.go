package config

import(
	"os"
	"github.com/golobby/dotenv"
)

type Config struct {
	Redis struct {
		Addr     string `env:"REDIS_ADDR"`
		Password string `env:"REDIS_PASSWORD"`
	}
	Mongodb struct {
		Port     string `env:"MONGO_PORT"`
		User     string `env:"MONGO_USER"`
		Password string `env:"MONGO_PASSWORD"`
	}
	GameLogic struct {
		GRPC_URL     string `env:"GAMELOGIC_GRPC_URL"`
		GRPC_GATEWAY_URL string `env:"GAMELOGIC_GRPC_GATEWAY_URL"`
	}
	Auth struct {
		GRPC_URL     string `env:"AUTH_GRPC_URL"`
	}
}

func GetConfig() (*Config) {
	config := Config{}
	file, err := os.Open("app.env")
	err = dotenv.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}
	return &config
}