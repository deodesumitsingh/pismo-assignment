package res

type ErrResp struct {
	Message string `json:"message,omitempty"`
}

type Resp struct {
	Data  any      `json:"data,omitempty"`
	Error *ErrResp `json:"error,omitempty"`
}

func NewResp(data any, err error) Resp {
	var r Resp
	r.Data = data

	if err != nil {
		r.Error = &ErrResp{
			Message: err.Error(),
		}
	}

	return r
}
