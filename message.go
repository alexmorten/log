package log

//IsInTimeRange checks if the message is in the timerange
func (m *Message) IsInTimeRange(startTime, endTime int64) bool {
	return m.Timestamp >= startTime && m.Timestamp <= endTime
}
