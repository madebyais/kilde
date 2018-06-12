package kilde

import (
	"encoding/json"
	"testing"
)

var (
	env          string
	configschema interface{}
	jsondata     string
)

type schema struct {
	Hostaddr string
	Domain   string
}

func TestFillStructSchemaShouldReturnErrorIfFieldIsNotValid(t *testing.T) {
	env = `default`
	jsondata = `{
    "default": {
      "hostaddr": "sample",
      "Domain": "www.kilde.doe"
    }
  }`

	_ = json.Unmarshal([]byte(jsondata), &configschema)

	s := &schema{}

	err := fillStructSchema(env, configschema, s)
	expectedError := "no such field `hostaddr` in config file"
	if err.Error() != expectedError {
		t.Error(`Expected error= `+expectedError+`, but got`, err.Error())
	}
}
