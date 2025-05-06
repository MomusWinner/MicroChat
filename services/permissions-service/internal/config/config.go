package config

import "github.com/alecthomas/kong"

type Config struct {
	KeyCloakURL   string `help:"Keycloak URL"                 env:"APP_KEYCLOAK_URL"         required:"true"`
	DatabaseURL   string `help:"Database URL"                 env:"APP_DATABASE_URL"         required:"true"`
	CloakClientId string `help:"KeyCloak client id"           env:"APP_CLOAK_CLIENT_ID"      required:"true"`
	CloakRealm    string `help:"KeyCloak client id"           env:"APP_CLOAK_REALM"          required:"true"`
	CloakSecret   string `help:"KeyCloak client id"           env:"APP_CLOAK_SECRET"         required:"true"`

	// GameServicePort string `help:"Port to listen on"                    env:"GAME_SERVICE_PORT"                 default:"7077"`
	// GRPCPort        string `help:"Port to listen on"                    env:"GRPC_PORT"         required:"true" default:"9090"`
	// LogLevel        string `help:"Log level (debug, info, warn, error)" env:"LOG_LEVEL"                         default:"info"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	parser, err := kong.New(cfg)
	if err != nil {
		return nil, err
	}

	// Parse command-line flags, environment variables, and config file
	_, err = parser.Parse(nil)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
