package config

type Config struct {
	Env                string
	Port               string
	UserTable          string
	AWS_ACCESS_KEY     string
	AWS_SECRET_KEY_ID  string
	AWS_DEFAULT_REGION string
	DatabaseName       string
	DatabaseUser       string
	DatabaseHost       string
	DatabasePort       int
	DatabaseRegion     string
	Debugging          bool
	Testing            bool
}
