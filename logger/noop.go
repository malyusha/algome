package logger

type noopLogger struct {
}

func (n *noopLogger) Debug(msg string, args ...any) {}
func (n *noopLogger) Info(msg string, args ...any)  {}
func (n *noopLogger) Error(msg string, args ...any) {}
func (n *noopLogger) Fatal(msg string, args ...any) {}
