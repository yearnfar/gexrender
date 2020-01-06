package util

import (
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

// Stringify json 序列化
func Stringify(v interface{}) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(v)
	if err != nil {
		log.Panic(err)
	}

	return string(data)
}
