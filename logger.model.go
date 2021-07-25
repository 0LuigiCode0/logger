package logger

import "os"

var Log Logerer

type Logerer interface {
	Info(args ...interface{})
	Service(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Servicef(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

//logger модель логгера
type logger struct {
	logs map[levelKey]*level
}

//level модель каждого уровня логгера
type level struct {
	format     string
	formatFile string
	file       *os.File
}

type levelKey string

const (
	LInfo    levelKey = "INFO"
	LService levelKey = "SERVICE"
	LWarning levelKey = "WARNING"
	LError   levelKey = "ERROR"
	LFatal   levelKey = "FATAL"
)

const (
	formatI = `{{font:b}}{{color:f}}{{color:b}}    {{level}} {{color:rb}} {{time:2006/01/02--15:04:05}} >>> WHERE
	{{file:1}}
MSG {{r}}"{{message}}"`
	formatS = `{{font:b}}{{color:b}}{{color:f}} {{level}} {{color:rb}} {{time:2006/01/02--15:04:05}} >>> WHERE
	{{file:1}}
MSG {{r}"{{message}}"`
	formatW = `{{font:b}}{{color:b}}{{color:f}} {{level}} {{color:rb}} {{time:2006/01/02--15:04:05}} >>> WHERE
	{{file:1}}
MSG {{r}}"{{message}}"`
	formatE = `{{font:b}}{{color:b}}{{color:f}}   {{level}} {{color:rb}} {{time:2006/01/02--15:04:05}} >>> WHERE
	{{file}}
MSG {{r}}"{{message}}"`
	formatF = `{{font:b}}{{color:b}}{{color:f}}   {{level}} {{color:rb}} {{time:2006/01/02--15:04:05}} >>> WHERE
	{{file}}
MSG {{r}}"{{message}}"`

	formatFI = `   {{level}} {{time:2006/01/02--15:04:05}} >>> WHERE {{file:1}} MSG {{message}}`
	formatFS = `{{level}} {{time:2006/01/02--15:04:05}} >>> WHERE {{file:1}} MSG {{message}}`
	formatFW = `{{level}} {{time:2006/01/02--15:04:05}} >>> WHERE {{file:1}} MSG {{message}}`
	formatFE = `  {{level}} {{time:2006/01/02--15:04:05}} >>> WHERE {{file}} MSG {{message}}`
	formatFF = `  {{level}} {{time:2006/01/02--15:04:05}} >>> WHERE {{file}} MSG {{message}}`
)
