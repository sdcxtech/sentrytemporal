package sentrytemporal

import "github.com/mitchellh/mapstructure"

func Struct2Map(st interface{}) (m map[string]interface{}, err error) {
	err = mapstructure.Decode(st, &m)
	return
}

func MustStruct2Map(st interface{}) map[string]interface{} {
	m, err := Struct2Map(st)
	if err != nil {
		panic(err)
	}

	return m
}
