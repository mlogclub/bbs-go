package simple

import "os"

type LogWriter struct {
	file *os.File
}

func NewLogWriter(logPath string) (*LogWriter, error) {
	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &LogWriter{file: file}, nil
}

func (w *LogWriter) Write(b []byte) (n int, err error) {
	_, _ = os.Stdout.Write(b)
	return w.file.Write(b)
}
