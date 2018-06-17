package kilde

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Kilde is an easy to use configuration library
type Kilde struct {
	configtype string
	filepath   string
	env        string

	filedata []byte

	schema       interface{}
	configschema interface{}
}

var (
	defaultconf *Kilde
)

// New is used to initiate kilde
func New() Interface {
	return &Kilde{
		env: `default`,
	}
}

// SetSchema is used to set configuration schema
func (k *Kilde) SetSchema(i interface{}) {
	k.schema = i
}

// GetSchema is used to get current defined schema
func (k *Kilde) GetSchema() interface{} {
	return k.schema
}

// GetConfigSchema is used to get default config schema
func (k *Kilde) GetConfigSchema() interface{} {
	return k.configschema
}

// SetConfigType is used to set configuration file type, e.g. json, yml, yaml, toml
func (k *Kilde) SetConfigType(t string) {
	k.configtype = t
}

// GetConfigType is used to get current config type
func (k *Kilde) GetConfigType() string {
	return k.configtype
}

// SetFilePath is used to set the configuration file path, hence kilde can locate the configuration file
func (k *Kilde) SetFilePath(filepath string) {
	k.filepath = filepath
}

// GetFilePath is used to get current defined file path
func (k *Kilde) GetFilePath() string {
	return k.filepath
}

// SetEnv is used to set application environment that want to be used
func (k *Kilde) SetEnv(env string) {
	k.env = env
}

// GetEnv is used to get current defined environment
func (k *Kilde) GetEnv() string {
	return k.env
}

// Read is used to initiate read and load configuration file
func (k *Kilde) Read() error {
	filedata, err := ioutil.ReadFile(k.filepath)
	if err != nil {
		return err
	}

	k.filedata = filedata

	var errParser error

	switch k.configtype {
	default:
		errParser = errors.New("cannot find any config type, please define one of these: json, yml, or yaml")
	case "json":
		errParser = k.json()
	case "yml", "yaml":
		errParser = k.yaml()
	}

	return errParser
}

func (k *Kilde) json() error {
	err := json.Unmarshal(k.filedata, &k.configschema)
	if err != nil {
		return err
	}

	return fill(k)
}

func (k *Kilde) yaml() error {
	var i map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(k.filedata), &i)
	if err != nil {
		return err
	}

	var temp = make(map[string]map[interface{}]interface{})
	var temp2 = make(map[string]map[string]interface{})
	var final = make(map[string]interface{})

	for key, val := range i {
		newKey := fmt.Sprintf(`%s`, key)
		temp[newKey] = val.(map[interface{}]interface{})
		temp2[newKey] = make(map[string]interface{})

		for i, v := range temp[newKey] {
			newKey2 := fmt.Sprintf(`%s`, i)
			temp2[newKey][newKey2] = v
		}

		final[newKey] = temp2[newKey]
	}

	k.configschema = final

	return fill(k)
}
