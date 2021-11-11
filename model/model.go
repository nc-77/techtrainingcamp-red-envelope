package model

type Envelope struct {
	EnvelopeId string `json:"envelope_id"`
	Value      int64  `json:"value"`
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

func (e *Envelope) ToMap() map[string]interface{} {
	ret := map[string]interface{}{
		"envelope_id": e.EnvelopeId,
		"value":       e.Value,
		"opened":      e.Opened,
		"snatch_time": e.SnatchTime,
		"user_id":     e.UserId,
	}
	return ret
}
