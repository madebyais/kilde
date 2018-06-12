package kilde

// Interface is an interface for Kilde
type Interface interface {
	SetSchema(i interface{})
	SetConfigType(t string)
	SetFilePath(filepath string)
	SetEnv(env string)

	GetSchema() interface{}
	GetConfigSchema() interface{}
	GetConfigType() string
	GetFilePath() string
	GetEnv() string

	Read() error
}
