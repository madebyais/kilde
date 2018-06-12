package main

import (
	"fmt"

	"gitlab.com/monologid/kilde"
)

type configschema struct {
	Host string `json:"Host" yaml:"Host"`
	Port string `json:"Port" yaml:"Port"`
}

func main() {
	schema := &configschema{}

	k := kilde.New()
	k.SetSchema(schema)
	k.SetConfigType(`json`)
	k.SetFilePath(`./example.json`)
	k.SetEnv(`development`)

	err := k.Read()
	if err != nil {
		panic(err)
	}

	fmt.Println(`--------- JSON ---------`)
	fmt.Println(`Host=`, schema.Host)
	fmt.Println(`Port=`, schema.Port)

	k.SetSchema(schema)
	k.SetConfigType(`yaml`)
	k.SetFilePath(`./example.yaml`)
	k.SetEnv(`staging`)

	err = k.Read()
	if err != nil {
		panic(err)
	}

	fmt.Println(`--------- YAML ---------`)
	fmt.Println(`Host=`, schema.Host)
	fmt.Println(`Port=`, schema.Port)
}
