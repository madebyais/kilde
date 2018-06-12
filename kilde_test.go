package kilde

import (
	"io/ioutil"
	"os"
	"testing"
)

type mockSchema struct {
	Hostaddr string  `json:"Hostaddr" yaml:"Hostaddr"`
	Port     float64 `json:"Port" yaml:"port"`
}

type mockSchemaYaml struct {
	Hostaddr string `yaml:"Hostaddr"`
	Port     int    `yaml:"port"`
}

var filepath string

func tearUp() {
	filepath = `./mockfile.json`

	jsondata := `{
		"default": {
			"Hostaddr": "127.0.0.1",
			"Port": 2000
		},
		"staging": {
			"Hostaddr": "staging.kilde.id"
		}
	}`

	err := ioutil.WriteFile(filepath, []byte(jsondata), 0644)
	if err != nil {
		panic(err)
	}
}

func tearUpInvalidJSON() {
	filepath = `./mockfile.json`

	jsondata := `{
		"default": {
			"hostaddr": "127.0.0.1"
		},
		"staging": {
			"hostaddr": "staging.kilde.id"
		}
	`

	err := ioutil.WriteFile(filepath, []byte(jsondata), 0644)
	if err != nil {
		panic(err)
	}
}

func tearUpIncorrectJSONField() {
	filepath = `./mockfile.json`

	jsondata := `{
		"default": {
			"hostaddr": "127.0.0.1"
		},
		"staging": {
			"hostaddr": "staging.kilde.id"
		}
	}`

	err := ioutil.WriteFile(filepath, []byte(jsondata), 0644)
	if err != nil {
		panic(err)
	}
}

var ymldata string = `
default:
  Hostaddr: "localhost"
  Port: 2000

staging:
  Hostaddr: "staging.kilde.id"
  Port: 9999
`

func tearUpYaml() {
	filepath = `./mockfile.yml`

	err := ioutil.WriteFile(filepath, []byte(ymldata), 0644)
	if err != nil {
		panic(err)
	}
}

var invalidymldata string = `
	default:
	  Hostaddr: "localhost"
	  Port: 2000

	staging:
	  Hostaddr: "staging.kilde.id"
	  Port: 9999
`

func tearUpInvalidYaml() {
	filepath = `./mockfile.yml`

	err := ioutil.WriteFile(filepath, []byte(invalidymldata), 0644)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	_ = os.Remove(filepath)
}

func TestSetSchemaShouldReturnCorrectSchemaAsDefined(t *testing.T) {
	k := New()
	s := &mockSchema{}
	k.SetSchema(s)

	if k.GetSchema().(*mockSchema) != s {
		t.Error(`Expected`, s, `but got`, k.GetSchema())
	}
}

func TestSetConfigTypeShouldReturnCorrectConfigTypeAsDefined(t *testing.T) {
	k := New()

	configtype := `json`
	k.SetConfigType(configtype)

	if k.GetConfigType() != configtype {
		t.Error(`Expected`, configtype, `but got`, k.GetConfigType())
	}

	configtype = `yaml`
	k.SetConfigType(configtype)

	if k.GetConfigType() != configtype {
		t.Error(`Expected`, configtype, `but got`, k.GetConfigType())
	}

	configtype = `yml`
	k.SetConfigType(configtype)

	if k.GetConfigType() != configtype {
		t.Error(`Expected`, configtype, `but got`, k.GetConfigType())
	}

	configtype = `toml`
	k.SetConfigType(configtype)

	if k.GetConfigType() != configtype {
		t.Error(`Expected`, configtype, `but got`, k.GetConfigType())
	}
}

func TestSetFilePathShouldReturnCorrectFilePathAsDefined(t *testing.T) {
	k := New()

	filepath = "/users/monologid/config.json"
	k.SetFilePath(filepath)

	if k.GetFilePath() != filepath {
		t.Error(`Expected`, filepath, `but got`, k.GetFilePath())
	}
}

func TestSetEnvShouldReturnCorrectEnvAsDefined(t *testing.T) {
	k := New()

	env := "development"
	k.SetEnv(env)

	if k.GetEnv() != env {
		t.Error(`Expected`, env, `but got`, k.GetEnv())
	}

	env = "staging"
	k.SetEnv(env)

	if k.GetEnv() != env {
		t.Error(`Expected`, env, `but got`, k.GetEnv())
	}

	env = "production"
	k.SetEnv(env)

	if k.GetEnv() != env {
		t.Error(`Expected`, env, `but got`, k.GetEnv())
	}
}

func TestReadShouldReturnErrorIfFileNotFound(t *testing.T) {
	k := New()
	k.SetFilePath("/users/monologid/mockfile.json")

	err := k.Read()

	if err == nil {
		t.Error(`Expected error but got no error`)
	}
}

func TestReadShouldReturnErrorIfConfigTypeIsNotSpecified(t *testing.T) {
	tearUp()

	s := &mockSchema{}

	k := New()
	k.SetSchema(s)
	k.SetFilePath(filepath)
	k.SetEnv(`staging`)

	err := k.Read()

	if err == nil {
		t.Error(`Expected error but got no error`)
	}

	tearDown()
}

func TestReadShouldReturnNoErrorIfFileIsFound(t *testing.T) {
	tearUp()

	s := &mockSchema{}

	k := New()
	k.SetSchema(s)
	k.SetFilePath(filepath)
	k.SetConfigType(`json`)
	k.SetEnv(`staging`)

	err := k.Read()

	if err != nil && err.Error() != "cannot find any config type, please define one of these: json, yml, yaml or toml" {
		t.Error(`Expected no error but got`, err.Error())
	}

	tearDown()
}

func TestReadShouldReadJson(t *testing.T) {
	tearUp()

	s := &mockSchema{}

	k := New()
	k.SetSchema(s)
	k.SetFilePath(filepath)
	k.SetConfigType(`json`)
	k.SetEnv(`staging`)

	err := k.Read()
	if err != nil {
		t.Error(`Expected no error but got`, err.Error())
	}

	cfg := k.GetConfigSchema().(map[string]interface{})

	if cfg["default"] == nil {
		t.Error(`Expected read json but got nothing`)
	}

	tearDown()
}

func TestReadShouldReturnErrorWhenReadJsonAndEnvNotAvailable(t *testing.T) {
	tearUp()

	s := &mockSchema{}

	k := New()
	k.SetSchema(s)
	k.SetFilePath(filepath)
	k.SetConfigType(`json`)
	k.SetEnv(`development`)

	err := k.Read()
	if err == nil {
		t.Error(`Expected an error but got no error`)
	}

	tearDown()
}

func TestReadShouldReturnErrorWhenReadFalseJson(t *testing.T) {
	tearUpInvalidJSON()

	k := New()
	k.SetFilePath(filepath)
	k.SetConfigType(`json`)

	err := k.Read()
	if err == nil {
		t.Error(`Expected an error but got no error`)
	}

	tearDown()
}

func TestReadShouldReadJsonButReturnError(t *testing.T) {
	tearUpIncorrectJSONField()

	s := &mockSchema{}

	k := New()
	k.SetSchema(s)
	k.SetFilePath(filepath)
	k.SetConfigType(`json`)
	k.SetEnv(`staging`)

	err := k.Read()
	if err == nil {
		t.Error(`Expected an error but got no error`)
	}

	tearDown()
}

func TestShouldReadYaml(t *testing.T) {
	tearUpYaml()

	s := &mockSchemaYaml{}

	k := New()
	k.SetSchema(s)
	k.SetFilePath(filepath)
	k.SetConfigType(`yml`)
	k.SetEnv(`staging`)

	err := k.Read()
	if err != nil {
		t.Error(`Expected no error but got`, err.Error())
	}

	tearDown()
}

func TestShouldReturnErrorWhenReadInvalidYaml(t *testing.T) {
	tearUpInvalidYaml()

	s := &mockSchemaYaml{}

	k := New()
	k.SetSchema(s)
	k.SetFilePath(filepath)
	k.SetConfigType(`yml`)
	k.SetEnv(`staging`)

	err := k.Read()
	if err == nil {
		t.Error(`Expected an error but got no error`)
	}

	tearDown()
}
