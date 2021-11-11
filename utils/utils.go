package utils

import (
	"encoding/json"
	"os"
	"strings"

	"red_envelope/model"
)

func Max(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func Min(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

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
