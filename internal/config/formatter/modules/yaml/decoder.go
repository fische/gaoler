package yaml

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Decoder struct {
	Reader io.Reader
}

func (d *Decoder) Decode(obj interface{}) error {
	data, err := ioutil.ReadAll(d.Reader)
	if err != nil {
		return err
	} else if len(data) == 0 {
		return io.EOF
	}
	return yaml.Unmarshal(data, obj)
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}
