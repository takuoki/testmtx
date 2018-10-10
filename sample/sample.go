package sample

// Request is ...
type Request struct {
	NumKey    *int    `json:"num_key"`
	StringKey *string `json:"string_key"`
	BoolKey   *bool   `json:"bool_key"`
	ObjectKey struct {
		Key1 *int    `json:"key1"`
		Key2 *string `json:"key2"`
	} `json:"object_key"`
	ArrayKey []struct {
		Key3 *int    `json:"key3"`
		Key4 *string `json:"key4"`
	} `json:"array_key"`
}

// Expected is ...
type Expected struct {
	Status *string `json:"status"`
	Code   *int    `json:"code"`
}

// Sample is ...
func Sample(req *Request) *Expected {

	var status string
	var code int
	if req != nil && req.NumKey != nil {
		switch *req.NumKey {
		case 101:
			status = "success"
			code = 200
		case 102:
			status = "failure"
			code = 401
		case 103:
			status = "failure"
			code = 404
		default:
			panic("unexpected number")
		}
	}

	return &Expected{
		Status: &status,
		Code:   &code,
	}
}
