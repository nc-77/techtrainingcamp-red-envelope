package model

type Envelope struct {
	EnvelopeId uint   `json:"envelope_id"`
	Value      uint   `json:"value,omitempty"`
	Opened     bool   `json:"opened"`
	SnatchTime string `json:"snatch_time"`
	UserId     uint   `json:"-"`
}
