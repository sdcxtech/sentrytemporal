package sentrytemporal

import "github.com/mitchellh/mapstructure"

func struct2Map(st interface{}) (m map[string]interface{}, err error) {
	err = mapstructure.Decode(st, &m)
	return
}

func mustStruct2Map(st interface{}) map[string]interface{} {
	m, err := struct2Map(st)
	if err != nil {
		panic(err)
	}

	return m
}
