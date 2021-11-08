package utils

import (
	"encoding/json"
	"os"
	"red_envelope/model"
	"strings"
)

func GetEnv(env string, defaultVal string) (key string) {
	if key = os.Getenv(env); key == "" {
		key = defaultVal
	}
	return
}

func GetArgs(env string) []string {
	return strings.Split(env, ";")
}

func DecodeWallet(envelopes map[string]string) (wallet []*model.Envelope, err error) {

	wallet = make([]*model.Envelope, len(envelopes))
	index := 0
	for _, envelope := range envelopes {
		if err = json.Unmarshal([]byte(envelope), &wallet[index]); err != nil {
			return nil, err
		}
		index++
	}
	return
}
