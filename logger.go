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

//Logger модель логгера
type Logger struct {
	logs map[levelKey]*level
}

//level модель каждого уровня логгера
type level struct {
	level      levelKey
	color      int
	backColor  int
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

//InitLogger иницилезирует логгер
func InitLogger(pathFile string) *Logger {
	Log := Logger{
		logs: make(map[levelKey]*level),
	}
	logI := new(level)
	logI.level = LInfo
	logI.color = 251
	logI.backColor = 233
	logI.format = logI.initForm("%{color:b}%{color:f}%{color}    %{level} %{color:rf} %{time:2006/01/02--15:04:05} >>> WHERE \n\t%{file:1}\n MSG %{color:r}\"%{message}\"")
	logI.formatFile = logI.initForm("   %{level} %{time:2006/01/02--15:04:05} >>> WHERE %{file:1} MSG %{message}")

	logS := new(level)
	logS.level = LService
	logS.color = 86
	logS.backColor = 96
	logS.format = logS.initForm("%{color:f}%{color:b}%{color} %{level} %{color:rf} %{time:2006/01/02--15:04:05} >>> WHERE \n\t%{file:1}\n MSG %{color:r}\"%{message}\"")
	logS.formatFile = logS.initForm("%{level} %{time:2006/01/02--15:04:05} >>> WHERE %{file:1} MSG %{message}")

	logW := new(level)
	logW.level = LWarning
	logW.color = 226
	logW.backColor = 239
	logW.format = logW.initForm("%{color:f}%{color:b}%{color} %{level} %{color:rf} %{time:2006/01/02--15:04:05} >>> WHERE \n\t%{file:1}\n MSG %{color:r}\"%{message}\"")
	logW.formatFile = logW.initForm("%{level} %{time:2006/01/02--15:04:05} >>> WHERE %{file:1} MSG %{message}")

	logE := new(level)
	logE.level = LError
	logE.color = 9
	logE.backColor = 188
	logE.format = logE.initForm("%{color:f}%{color:b}%{color}   %{level} %{color:rf} %{time:2006/01/02--15:04:05} >>> WHERE \n\t%{file}\n MSG %{color:r}\"%{message}\"")
	logE.formatFile = logE.initForm("  %{level} %{time:2006/01/02--15:04:05} >>> WHERE %{file} MSG %{message}")

	logF := new(level)
	logF.level = LFatal
	logF.color = 128
	logF.backColor = 215
	logF.format = logF.initForm("%{color:f}%{color:b}%{color}   %{level} %{color:rf} %{time:2006/01/02--15:04:05} >>> WHERE \n\t%{file}\n MSG %{color:r}\"%{message}\"")
	logF.formatFile = logF.initForm("  %{level} %{time:2006/01/02--15:04:05} >>> WHERE %{file} MSG %{message}")

	Log.logs[LInfo] = logI
	Log.logs[LService] = logS
	Log.logs[LWarning] = logW
	Log.logs[LError] = logE
	Log.logs[LFatal] = logF

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

//Infof выводит форматируемое информационное сообщение
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logs["info"].log(format, args...)
}

//Servicef выводит форматируемое серверное сообщение
func (l *Logger) Servicef(format string, args ...interface{}) {
	l.logs["serv"].log(format, args...)
}

//Warningf выводит форматируемое предупреждающее сообщение
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.logs["warn"].log(format, args...)
}

//Errorf выводит форматируемое сообщение ошибки
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logs["err"].log(format, args...)
}

//Fatalf выводит форматируемое критическое сообщение и завершает программу
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

//Fatal выводит критическое сообщение и завершает программу
func (l *Logger) Fatal(args ...interface{}) {
	l.logs["fatal"].log("", args...)
	os.Exit(1)
}

//log оснавная функция вывода сообщения в консоль с учетом форматирования
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

//SetLevelFormatConsole изменение формата вывода в консоль по конкретному уровню
func (l *Logger) SetLevelFormatConsole(level levelKey, format string) {
	if lv, ok := l.logs[level]; ok {
		lv.format = format
	} else {
		l.Warning("Level is not found")
	}
}

//SetLevelFormatFile изменение формата вывода в фаил по конкретному уровню
func (l *Logger) SetLevelFormatFile(level levelKey, format string) {
	if lv, ok := l.logs[level]; ok {
		lv.formatFile = format
	} else {
		l.Warning("Level is not found")
	}
}

//SetLevelFile изменение файла вывода по конкретному уровню
func (l *Logger) SetLevelFile(level levelKey, file *os.File) {
	if lv, ok := l.logs[level]; ok {
		lv.file = file
	} else {
		l.Warning("Level is not found")
	}
}

//SetLevelColor изменение цвета сообщения по конкретному уровню
func (l *Logger) SetLevelColor(level levelKey, color int) {
	if lv, ok := l.logs[level]; ok {
		lv.color = color
	} else {
		l.Warning("Level is not found")
	}
}

//SetLevelBackColor изменение цвета фона сообщения по конкретному уровню
func (l *Logger) SetLevelBackColor(level levelKey, color int) {
	if lv, ok := l.logs[level]; ok {
		lv.backColor = color
	} else {
		l.Warning("Level is not found")
	}
}

//SetFormatConsole изменение формата вывода в консоль для всех уровней
func (l *Logger) SetFormatConsole(format string) {
	for _, v := range l.logs {
		v.format = format
	}
}

//SetFormatFile изменение формата вывода в фаил  для всех уровней
func (l *Logger) SetFormatFile(format string) {
	for _, v := range l.logs {
		v.formatFile = format
	}
}

//SetFile изменение файла вывода  для всех уровней
func (l *Logger) SetFile(file *os.File) {
	for _, v := range l.logs {
		v.file = file
	}
}

//initForm парсит основное форматирование уровня, заполняя управляющими байтами
func (lv *level) initForm(format string) string {
	color := "\033[38;5;" + fmt.Sprint(lv.color) + "m"
	colorF := "\033[48;5;" + fmt.Sprint(lv.backColor) + "m"
	format = strings.ReplaceAll(format, "%{message}", "%v")
	format = strings.ReplaceAll(format, "%{level}", string(lv.level))
	format = strings.ReplaceAll(format, "%{color}", color)
	format = strings.ReplaceAll(format, "%{color:f}", colorF)
	format = strings.ReplaceAll(format, "%{color:rf}", "\033[49m")
	format = strings.ReplaceAll(format, "%{color:rb}", "\033[22m")
	format = strings.ReplaceAll(format, "%{color:r}", "\033[0m")
	format = strings.ReplaceAll(format, "%{color:b}", "\033[1m")
	return format
}

//parser парсит форматирование сообщения, заполняя необходимыми данными
func (lv *level) parser(format string) string {
	regex := regexp.MustCompile("%{time[^{}]+")
	dateFind := regex.FindString(format)
	regex = regexp.MustCompile("%{file[^{}]+")
	fileFind := regex.FindString(format)
	dateForm := strings.TrimPrefix(dateFind, "%{time:")
	fileCalldeph, _ := strconv.Atoi(strings.TrimPrefix(fileFind, "%{file:"))
	t := time.Now()
	tForm := t.Format("01 02 2006:04:05.000000")
	if dateForm != "" {
		tForm = t.Format(dateForm)
	}
	if dateFind == "" {
		dateFind = "%{time"
	}
	if fileFind == "" {
		fileFind = "%{file"
	}
	pc := make([]uintptr, 6)
	runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc)
	trace := []string{}
	for f, b := frames.Next(); b; f, b = frames.Next() {
		frame := ""
		arrHistrF := strings.Split(f.File, "/")
		if fileCalldeph > 0 && fileCalldeph < len(arrHistrF) {
			for i := 1; i <= fileCalldeph; i++ {
				frame += "/" + arrHistrF[len(arrHistrF)-i]
			}
		} else {
			frame = f.File
		}
		funcName := regexp.MustCompile(`[^\.]*$`).FindString(f.Function)
		trace = append(trace, fmt.Sprintf("%v:%v (%v)", frame, fmt.Sprint(f.Line), funcName))
		if funcName == "main" {
			break
		}
	}
	format = strings.ReplaceAll(format, dateFind+"}", tForm)
	format = strings.ReplaceAll(format, fileFind+"}", strings.Join(trace, "\n\t⮤ "))
	return format
}
