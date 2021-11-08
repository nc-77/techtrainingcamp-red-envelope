package utils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"red_envelope/model"
	"testing"
)

func TestDecodeWallet(t *testing.T) {
	envelopes := []*model.Envelope{
		{},
		{EnvelopeId: "1", Opened: true},
		{EnvelopeId: "2", Value: 1},
	}
	input := make(map[string]string)
	for _, envelope := range envelopes {
		data, _ := json.Marshal(envelope)
		input[envelope.EnvelopeId] = string(data)
	}
	got, err := DecodeWallet(input)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, envelopes, got)
}
