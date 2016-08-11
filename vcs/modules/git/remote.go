package git

//Remote represents a git remote
type Remote struct {
	Name string
	URL  string
	Type byte
}

const (
	fetch = 0x01
	push  = 0x02
)

//IsFetch checks if remote is of type fetch
func (r Remote) IsFetch() bool {
	return (r.Type & fetch) != 0
}

//IsPush checks if remote is of type push
func (r Remote) IsPush() bool {
	return (r.Type & push) != 0
}

//GetURL returns the remote URL
func (r Remote) GetURL() string {
	return r.URL
}
