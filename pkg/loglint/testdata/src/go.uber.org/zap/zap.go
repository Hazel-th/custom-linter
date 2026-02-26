package zap

// Это минимальный stub zap, который используется только в analysistest-фикстурах.
// Он нужен, чтобы проходила проверка типов без подключения реального zap в testdata.
type Field struct{}

type Logger struct{}

type SugaredLogger struct{}

func L() *Logger { return &Logger{} }

func (l *Logger) Debug(msg string, fields ...Field)  {}
func (l *Logger) Info(msg string, fields ...Field)   {}
func (l *Logger) Warn(msg string, fields ...Field)   {}
func (l *Logger) Error(msg string, fields ...Field)  {}
func (l *Logger) DPanic(msg string, fields ...Field) {}
func (l *Logger) Panic(msg string, fields ...Field)  {}
func (l *Logger) Fatal(msg string, fields ...Field)  {}

func (s *SugaredLogger) Debug(args ...any)  {}
func (s *SugaredLogger) Info(args ...any)   {}
func (s *SugaredLogger) Warn(args ...any)   {}
func (s *SugaredLogger) Error(args ...any)  {}
func (s *SugaredLogger) DPanic(args ...any) {}
func (s *SugaredLogger) Panic(args ...any)  {}
func (s *SugaredLogger) Fatal(args ...any)  {}

func (s *SugaredLogger) Debugw(msg string, keysAndValues ...any)  {}
func (s *SugaredLogger) Infow(msg string, keysAndValues ...any)   {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...any)   {}
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...any)  {}
func (s *SugaredLogger) DPanicw(msg string, keysAndValues ...any) {}
func (s *SugaredLogger) Panicw(msg string, keysAndValues ...any)  {}
func (s *SugaredLogger) Fatalw(msg string, keysAndValues ...any)  {}
