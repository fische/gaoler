package json

import "encoding/json"

type Encoder struct {
	*json.Encoder
}

func (e *Encoder) PrettyEncode(obj interface{}) error {
	e.SetIndent("", "\t")
	defer e.SetIndent("", "")
	return e.Encode(obj)
}
