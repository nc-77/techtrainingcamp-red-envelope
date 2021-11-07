package model

type Envelope struct {
	EnvelopeId string `json:"envelope_id"`
	Value      int64  `json:"value,omitempty"`
	Opened     bool   `json:"opened"`
	SnatchTime int64  `json:"snatch_time"`
	UserId     string `json:"user_id"`
}

type RespEnvelope struct {
	EnvelopeId string `json:"envelope_id"`
	Value      int64  `json:"value,omitempty"`
	Opened     bool   `json:"opened"`
	SnatchTime int64  `json:"snatch_time"`
}
