package utils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"red_envelope/model"
	"testing"
)

func TestDecodeWallet(t *testing.T) {
	envelopes := []*model.Envelope{
		{EnvelopeId: "1", Opened: true, Value: 1},
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
	for i := range envelopes {
		assert.Equal(t, *envelopes[i], *got[i])
	}

}
