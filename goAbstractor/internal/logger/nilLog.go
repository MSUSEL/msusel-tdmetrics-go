package logger

type nilLog struct{}

func (nilLog) Log(format string, args ...any) {}
