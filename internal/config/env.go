package config

const (
	devFile  = "development.json"
	prodFile = "production.json"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func SetEnv(envType string) string {
	switch envType {
	case EnvDev:
		return devFile
	case EnvProd:
		return prodFile
	default:
		return devFile
	}
}
