package kilde

import (
	"errors"
	"fmt"
	"reflect"
)

func fill(k *Kilde) error {
	err := fillStructSchema(`default`, k.configschema, k.schema)
	if err != nil {
		return err
	}

	err = fillStructSchema(k.env, k.configschema, k.schema)
	if err != nil {
		return err
	}

	return nil
}

func fillStructSchema(env string, configschema interface{}, schema interface{}) error {
	cfg := configschema.(map[string]interface{})

	if _, ok := cfg[env]; !ok {
		fmt.Println(fmt.Sprintf(`[madebyais/kilde] %s not found, set env to default`, env))
		env = `default`
	}

	structVal := reflect.ValueOf(schema).Elem()
	for key, value := range cfg[env].(map[string]interface{}) {
		structFieldValue := structVal.FieldByName(key)

		if !structFieldValue.IsValid() {
			return errors.New("no such field `" + key + "` in config file")
		}

		if !structFieldValue.CanSet() {
			return errors.New("cannot set `" + key + "` field value")
		}

		val := reflect.ValueOf(value)
		if structFieldValue.Type() != val.Type() {
			return errors.New("provided value type didn't match `" + key + "` obj field type")
		}

		structFieldValue.Set(val)
	}

	return nil
}
