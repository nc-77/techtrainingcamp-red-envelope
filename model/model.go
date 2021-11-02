package model

type Envelope struct {
	EnvelopeId uint64 `json:"envelope_id"`
	Value      uint64 `json:"value,omitempty"`
	Opened     bool   `json:"opened"`
	SnatchTime string `json:"snatch_time"`
	UserId     uint64 `json:"-"`
}
