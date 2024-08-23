package config

type Config struct {
	DatabaseUser string
	DatabasePsw  string
	DatabaseURL  string
	DatabasePORT string
	DatabaseName string
}

func LoadConfig() (*Config, error) {

	return &Config{
		DatabaseUser: "root",
		DatabasePsw:  "admin",
		DatabaseURL:  "127.0.0.1",
		DatabasePORT: "3306",
		DatabaseName: "dbtodo",
	}, nil

}
