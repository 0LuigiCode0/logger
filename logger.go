package logger

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	black = iota
	red
	green
	yellow
	blue
	purple
	cyan
	white
)

const (
	fBlack = iota
	fRed
	fGreen
	fYellow
	fBlue
	fPurple
	fCyan
	fWhite
)

//Logger класс логгера
type Logger struct {
	logs map[string]*level
}

type level struct {
	level      string
	color      int
	backColor  int
	format     string
	formatFile string
	file       *os.File
}

//InitLogger иницилезирует логгер
func InitLogger(pathFile string) *Logger {
	Log := Logger{
		logs: make(map[string]*level),
	}
	logI := new(level)
	logI.level = "INFO"
	logI.color = 251
	logI.backColor = 233
	logI.format = logI.initForm("%{color:b}%{color:f}%{color}    %{level} %{color:rf} %{time:2006/01/02--15:04:05} >>> WHERE [ %{file:1} ] MSG %{color:r}\"%{message}\"")
	logI.formatFile = logI.initForm("   %{level} %{time:2006/01/02--15:04:05} >>> WHERE %{file:1} MSG %{message}")

	logS := new(level)
	logS.level = "SERVICE"
	logS.color = 86
	logS.backColor = 96
	logS.format = logS.initForm("%{color:f}%{color:b}%{color} %{level} %{color:rf} %{time:2006/01/02--15:04:05} >>> WHERE [ %{file:1} ] MSG %{color:r}\"%{message}\"")
	logS.formatFile = logS.initForm("%{level} %{time:2006/01/02--15:04:05} >>> WHERE %{file:1} MSG %{message}")

	logW := new(level)
	logW.level = "WARNING"
	logW.color = 226
	logW.backColor = 239
	logW.format = logW.initForm("%{color:f}%{color:b}%{color} %{level} %{color:rf} %{time:2006/01/02--15:04:05} >>> WHERE [ %{file:1} ] MSG %{color:r}\"%{message}\"")
	logW.formatFile = logW.initForm("%{level} %{time:2006/01/02--15:04:05} >>> WHERE %{file:1} MSG %{message}")

	logE := new(level)
	logE.level = "ERROR"
	logE.color = 9
	logE.backColor = 188
	logE.format = logE.initForm("%{color:f}%{color:b}%{color}   %{level} %{color:rf} %{time:2006/01/02--15:04:05} >>> WHERE [ %{file} ] MSG %{color:r}\"%{message}\"")
	logE.formatFile = logE.initForm("  %{level} %{time:2006/01/02--15:04:05} >>> WHERE %{file} MSG %{message}")

	logF := new(level)
	logF.level = "FATAL"
	logF.color = 128
	logF.backColor = 215
	logF.format = logF.initForm("%{color:f}%{color:b}%{color}   %{level} %{color:rf} %{time:2006/01/02--15:04:05} >>> WHERE [ %{file} ] MSG %{color:r}\"%{message}\"")
	logF.formatFile = logF.initForm("  %{level} %{time:2006/01/02--15:04:05} >>> WHERE %{file} MSG %{message}")

	Log.logs["info"] = logI
	Log.logs["warn"] = logW
	Log.logs["err"] = logE
	Log.logs["fatal"] = logF
	Log.logs["serv"] = logS

	if pathFile != "" {
		outFile, err := os.OpenFile(pathFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			Log.Info("Log canot open")
		} else {
			logI.file = outFile
			logW.file = outFile
			logE.file = outFile
			logF.file = outFile
			logS.file = outFile
		}
	}

	return &Log
}

//Infof выводит информационное сообщение
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logs["info"].log(format, args...)
}

//Servicef выводит серверное сообщение
func (l *Logger) Servicef(format string, args ...interface{}) {
	l.logs["serv"].log(format, args...)
}

//Warningf выводит предупреждающее сообщение
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.logs["warn"].log(format, args...)
}

//Errorf выводит сообщение ошибки
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logs["err"].log(format, args...)
}

//Fatalf выводит критическое сообщение
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logs["fatal"].log(format, args...)
	os.Exit(1)
}

//Info выводит информационное сообщение
func (l *Logger) Info(args ...interface{}) {
	l.logs["info"].log("", args...)
}

//Service выводит серверное сообщение
func (l *Logger) Service(args ...interface{}) {
	l.logs["serv"].log("", args...)
}

//Warning выводит предупреждающее сообщение
func (l *Logger) Warning(args ...interface{}) {
	l.logs["warn"].log("", args...)
}

//Error выводит сообщение ошибки
func (l *Logger) Error(args ...interface{}) {
	l.logs["err"].log("", args...)
}

//Fatal выводит критическое сообщение
func (l *Logger) Fatal(args ...interface{}) {
	l.logs["fatal"].log("", args...)
	os.Exit(1)
}

func (lv *level) log(format string, args ...interface{}) {
	var message string
	if strings.TrimSpace(format) != "" {
		message = fmt.Sprintf(format, args...)
	} else {
		for _, v := range args {
			if message == "" {
				message = fmt.Sprint(v)
			} else {
				message = strings.Join([]string{message, fmt.Sprint(v)}, " ")
			}
		}
	}
	formatlevel := lv.parser(lv.format)
	fmt.Printf(formatlevel+"\n", message)
	if lv.file != nil {
		formatFile := lv.parser(lv.formatFile)
		_, err := lv.file.WriteString(fmt.Sprintf(formatFile+"\n", message))
		if err != nil {
			fmt.Printf("\n"+formatlevel, "Log writer failed")
		}
	}
}

//SetLevelFormatConsole изменение формата вывода в консоль
func (l *Logger) SetLevelFormatConsole(level, format string) {
	if lv, ok := l.logs[level]; ok {
		lv.format = format
	} else {
		l.Warning("Level is not found")
	}
}

//SetLevelFormatFile изменение формата вывода в фаил
func (l *Logger) SetLevelFormatFile(level, format string) {
	if lv, ok := l.logs[level]; ok {
		lv.formatFile = format
	} else {
		l.Warning("Level is not found")
	}
}

//SetLevelFile изменение файла вывода
func (l *Logger) SetLevelFile(level string, file *os.File) {
	if lv, ok := l.logs[level]; ok {
		lv.file = file
	} else {
		l.Warning("Level is not found")
	}
}

//SetLevelColor изменение цвета сообщения
func (l *Logger) SetLevelColor(level string, color int) {
	if lv, ok := l.logs[level]; ok {
		lv.color = color
	} else {
		l.Warning("Level is not found")
	}
}

//SetLevelBackColor изменение цвета фона сообщения
func (l *Logger) SetLevelBackColor(level string, color int) {
	if lv, ok := l.logs[level]; ok {
		lv.backColor = color
	} else {
		l.Warning("Level is not found")
	}
}

//SetFormatConsole изменение формата вывода в консоль
func (l *Logger) SetFormatConsole(format string) {
	for _, v := range l.logs {
		v.format = format
	}
}

//SetFormatFile изменение формата вывода в фаил
func (l *Logger) SetFormatFile(format string) {
	for _, v := range l.logs {
		v.formatFile = format
	}
}

//SetFile изменение файла вывода
func (l *Logger) SetFile(file *os.File) {
	for _, v := range l.logs {
		v.file = file
	}
}

func (lv *level) initForm(format string) string {
	color := "\x1b[38;5;" + fmt.Sprint(lv.color) + "m"
	colorF := "\x1b[48;5;" + fmt.Sprint(lv.backColor) + "m"
	format = strings.ReplaceAll(format, "%{message}", "%v")
	format = strings.ReplaceAll(format, "%{level}", lv.level)
	format = strings.ReplaceAll(format, "%{color}", color)
	format = strings.ReplaceAll(format, "%{color:f}", colorF)
	format = strings.ReplaceAll(format, "%{color:rf}", "\x1b[49m")
	format = strings.ReplaceAll(format, "%{color:rb}", "\x1b[22m")
	format = strings.ReplaceAll(format, "%{color:r}", "\x1b[0m")
	format = strings.ReplaceAll(format, "%{color:b}", "\x1b[1m")
	return format
}

func (lv *level) parser(format string) string {
	regex := regexp.MustCompile("%{time[^{}]+")
	dateFind := regex.FindString(format)
	regex = regexp.MustCompile("%{file[^{}]+")
	fileFind := regex.FindString(format)
	dateForm := strings.TrimPrefix(dateFind, "%{time:")
	fileCalldeph, _ := strconv.ParseInt(strings.TrimPrefix(fileFind, "%{file:"), 10, 64)
	t := time.Now()
	tForm := t.Format("01 02 2006:04:05.000000")
	if dateForm != "" {
		tForm = t.Format(dateForm)
	}
	if dateFind == "" {
		dateFind = "%{time"
	}
	_, file, line, ok := runtime.Caller(3)
	if ok {
		arrHistrF := strings.Split(file, "/")
		if fileCalldeph > 0 && fileCalldeph < int64(len(arrHistrF)) {
			file = ""
			for i := 1; int64(i) <= fileCalldeph; i++ {
				file = "/" + arrHistrF[len(arrHistrF)-i] + file
			}
		}
	}
	if fileFind == "" {
		fileFind = "%{file"
	}
	file += ":" + fmt.Sprint(line)
	format = strings.ReplaceAll(format, dateFind+"}", tForm)
	format = strings.ReplaceAll(format, fileFind+"}", file)
	return format
}
