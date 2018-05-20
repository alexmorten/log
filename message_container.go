package log

//MessageContainer is one of [PlainMessage, ServiceMessage, CompleteMessage]
type MessageContainer interface {
	GetLogMessage() *Message
}

//GetLogMessage ...
func (m *PlainMessage) GetLogMessage() *Message {
	return m.Message
}

//GetLogMessage ...
func (m *ServiceMessage) GetLogMessage() *Message {
	return m.Message
}

//GetLogMessage ...
func (m *CompleteMessage) GetLogMessage() *Message {
	return m.Message
}
