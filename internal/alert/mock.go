package alert

import "log"

// MockAlert .
type MockAlert struct{}

// Alert .
func (a *MockAlert) Alert(message Message) {
	log.Println("text : ", message.Text)
	log.Println("error : ", message.Error.Error())
	log.Println("trace : ", string(message.Trace))
}

// NewMockAlert .
func NewMockAlert() *MockAlert {
	return &MockAlert{}
}
