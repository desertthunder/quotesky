package log

import (
	"fmt"
	"time"
)

const (
	Reset string = "\033[0m"
	// Pine: rgb(49, 116, 143) as bg for DEBUG
	Debug     string = "\033[48;2;49;116;143m"
	DebugText string = "\033[38;2;49;116;143m"
	// Gold: rgb(246, 193, 119)
	Info     string = "\033[48;2;246;193;119m"
	InfoText string = "\033[38;2;246;193;119m"
	// Iris: rgb(196, 167, 231)
	Warn     string = "\033[48;2;196;167;231m"
	WarnText string = "\033[38;2;196;167;231m"
	// Love: rgb(235, 111, 146)
	Error     string = "\033[48;2;235;111;146m"
	ErrorText string = "\033[38;2;235;111;146m"
	// Rose (Dawn): rgb(215, 130, 126)
	Fatal     string = "\033[48;2;215;130;126m"
	FatalText string = "\033[38;2;215;130;126m"
	// Text: rgb(224, 222, 244)
	Text string = "\033[38;2;224;222;244m"
)

type LogLevel int32

type Logger struct {
	Level LogLevel
}

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String formats the log level.
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBU"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERRO"
	case FATAL:
		return "FATA"
	default:
		return ""
	}
}

// Tag sets the background color for the log level tag.
func (l LogLevel) Tag() string {
	t := l.String()
	switch l {
	case DEBUG:
		return fmt.Sprintf("%s %s %s ", Debug, t, Reset)
	case INFO:
		return fmt.Sprintf("%s %s %s ", Info, t, Reset)
	case WARN:
		return fmt.Sprintf("%s %s %s ", Warn, t, Reset)
	case ERROR:
		return fmt.Sprintf("%s %s %s ", Error, t, Reset)
	case FATAL:
		return fmt.Sprintf("%s %s %s ", Fatal, t, Reset)
	default:
		return "UNKNOWN"
	}
}

// Time formats the current (log) time and changes the
// color to match the tag background.
func (l LogLevel) Time() string {
	t := time.Now().Format(time.RFC822)
	switch l {
	case DEBUG:
		return fmt.Sprintf("%s[%s]%s ", DebugText, t, Reset)
	case INFO:
		return fmt.Sprintf("%s[%s]%s ", InfoText, t, Reset)
	case WARN:
		return fmt.Sprintf("%s[%s]%s ", WarnText, t, Reset)
	case ERROR:
		return fmt.Sprintf("%s[%s]%s ", ErrorText, t, Reset)
	case FATAL:
		return fmt.Sprintf("%s[%s]%s ", FatalText, t, Reset)
	default:
		return fmt.Sprintf("[%s] ", t)
	}
}

func (l Logger) Message(s string) string {
	return fmt.Sprintf("%s%s%s", Text, s, Reset)
}

func (l Logger) print(lv LogLevel, msg string) {
	if l.Level > lv {
		return
	}

	tag := lv.Tag()
	time := lv.Time()
	message := l.Message(msg)

	fmt.Print(tag, time)
	fmt.Println(message)
}

func (l Logger) Debug(msg string) {
	l.print(DEBUG, msg)
}

func (l Logger) Info(msg string) {
	l.print(INFO, msg)
}

func (l Logger) Warn(msg string) {
	l.print(WARN, msg)
}

func (l Logger) Error(msg string) {
	l.print(ERROR, msg)
}

func (l Logger) Fatal(msg string) {
	l.print(FATAL, msg)
}

func (l *Logger) SetLevel(lv LogLevel) {
	l.Level = lv
}

func NewLogger(level *LogLevel) *Logger {
	if level != nil {
		return &Logger{*level}
	} else {

		return &Logger{DEBUG}
	}
}

func DefaultLogger() *Logger {
	return NewLogger(nil)
}
