package server

import "encoding/json"

type RestResponse struct {
	Success bool
	Errors  []string
}

func NewRestResponse() *RestResponse {
	return &RestResponse{
		Success: true,
		Errors:  make([]string, 0),
	}
}

func (r *RestResponse) Error(err error) {
	r.Success = false
	r.Errors = append(r.Errors, err.Error())
}

func (r *RestResponse) Build() []byte {
	v, _ := json.Marshal(r)
	return v
}
