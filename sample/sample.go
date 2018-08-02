package sample

// Request is ...
type Request struct {
	IntKey    int    `json:"int_key"`
	StringKey string `json:"string_key"`
	BoolKey   bool   `json:"bool_key"`
	ObjectKey struct {
		Key1 int    `json:"key1"`
		Key2 string `json:"key2"`
	} `json:"object_key"`
	ArrayKey []struct {
		Key3 int    `json:"key3"`
		Key4 string `json:"key4"`
	} `json:"array_key"`
}

// Expected is ...
type Expected struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

// Sample is ...
func Sample(req *Request) *Expected {
	e := &Expected{}
	switch req.IntKey {
	case 101:
		e.Status = "success"
		e.Code = 200
	case 102:
		e.Status = "failure"
		e.Code = 401
	case 103:
		e.Status = "failure"
		e.Code = 404
	default:
		panic("unexpected number")
	}
	return e
}
