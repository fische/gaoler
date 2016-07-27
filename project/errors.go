package project

//ErrorMessage represents an error message
type ErrorMessage struct {
	Err     error
	Message string
	Fields  map[string]interface{}
}

//NewErrorMessage creates a new error message using given error
func NewErrorMessage(e error) *ErrorMessage {
	return &ErrorMessage{
		Err:    e,
		Fields: make(map[string]interface{}),
	}
}

//Error returns error message from `Err` field
func (m *ErrorMessage) Error() string {
	return m.Err.Error()
}

//WithFields merges given map to `Fields` map
func (m *ErrorMessage) WithFields(fields map[string]interface{}) *ErrorMessage {
	for k, v := range fields {
		m.Fields[k] = v
	}
	return m
}

//WithField sets `value` at `key` in `Fields` map
func (m *ErrorMessage) WithField(key string, value interface{}) *ErrorMessage {
	m.Fields[key] = value
	return m
}

//WithMessage sets `Message` field
func (m *ErrorMessage) WithMessage(msg string) *ErrorMessage {
	m.Message = msg
	return m
}
