package util

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/la0rg/test_tasks/vector_clock"
)

// vc_coder is intended to be used for encoding/decoding vector clocks
// to be able to share them with clients

// EncodeVc encodes vector clock to base64 representation
func EncodeVc(vc *vector_clock.VC) (string, error) {
	b, err := json.Marshal(vc)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(b), nil
}

// DecodeVc decodes base64 representation of vector clock
func DecodeVc(str string) (*vector_clock.VC, error) {
	bytes, err := b64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	var vc *vector_clock.VC
	err = json.Unmarshal(bytes, vc)
	if err != nil {
		return nil, err
	}
	return vc, nil
}
