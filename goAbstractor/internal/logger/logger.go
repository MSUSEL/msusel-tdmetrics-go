package logger

type Logger interface {
	Log(format string, args ...any)
}

func New(verbod bool) Logger {
	if verbod {
		return outLog{}
	}
	return nilLog{}
}
