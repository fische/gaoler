package yaml

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Decoder struct {
	r io.Reader
}

func (e *Decoder) Decode(obj interface{}) error {
	data, err := ioutil.ReadAll(e.r)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return io.EOF
	}
	return yaml.Unmarshal(data, obj)
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}
