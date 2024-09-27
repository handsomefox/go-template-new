package env

import (
	"os"
)

type Environment uint8

const (
	_ Environment = iota
	Local
	Testing
	Staging
	Development
	CI
	Production
)

const (
	Key = "ENVIRONMENT"
)

func Get() Environment {
	switch os.Getenv(Key) {
	case "local":
		return Local
	case "testing":
		return Testing
	case "staging":
		return Staging
	case "development", "dev":
		return Development
	case "ci", "workflows":
		return CI
	case "production", "prod":
		return Production
	default:
		return Local
	}
}
