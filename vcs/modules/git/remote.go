package git

type Remote struct {
	Name string
	URL  string
	Type byte
}

const (
	fetch = 1 << iota
	push
)

func (r Remote) IsFetch() bool {
	return (r.Type & fetch) != 0
}

func (r Remote) IsPush() bool {
	return (r.Type & push) != 0
}

func (r Remote) GetURL() string {
	return r.URL
}
