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

func init() {
	Log = InitLogger(nil)
}

//InitLogger иницилезирует логгер если pathFile = "", то фаил лога не создастся
func InitLogger(option Options) Logerer {
	return &logger{
		logs: map[levelKey]*level{
			LInfo:    newLevel(LInfo, 251, 233, formatI, formatFI, option),
			LService: newLevel(LService, 86, 96, formatS, formatFS, option),
			LWarning: newLevel(LWarning, 226, 239, formatW, formatFW, option),
			LError:   newLevel(LError, 9, 188, formatE, formatFE, option),
			LFatal:   newLevel(LFatal, 128, 215, formatF, formatFF, option),
		},
	}
}

//Infof выводит форматируемое информационное сообщение
func (l *logger) Infof(format string, args ...interface{}) {
	l.logs[LInfo].log(format, args...)
}

//Servicef выводит форматируемое серверное сообщение
func (l *logger) Servicef(format string, args ...interface{}) {
	l.logs[LService].log(format, args...)
}

//Warningf выводит форматируемое предупреждающее сообщение
func (l *logger) Warningf(format string, args ...interface{}) {
	l.logs[LWarning].log(format, args...)
}

//Errorf выводит форматируемое сообщение ошибки
func (l *logger) Errorf(format string, args ...interface{}) {
	l.logs[LError].log(format, args...)
}

//Fatalf выводит форматируемое критическое сообщение и завершает программу
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logs[LFatal].log(format, args...)
	os.Exit(1)
}

//Info выводит информационное сообщение
func (l *logger) Info(args ...interface{}) {
	l.logs[LInfo].log("", args...)
}

//Service выводит серверное сообщение
func (l *logger) Service(args ...interface{}) {
	l.logs[LService].log("", args...)
}

//Warning выводит предупреждающее сообщение
func (l *logger) Warning(args ...interface{}) {
	l.logs[LWarning].log("", args...)
}

//Error выводит сообщение ошибки
func (l *logger) Error(args ...interface{}) {
	l.logs[LError].log("", args...)
}

//Fatal выводит критическое сообщение и завершает программу
func (l *logger) Fatal(args ...interface{}) {
	l.logs[LFatal].log("", args...)
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

func newLevel(levelKey levelKey, color, backColor int, format, formatFile string, option Options) *level {
	l := &level{}
	if option != nil {
		if opt := option.get(levelKey); opt != nil {
			if opt.format != nil {
				format = *opt.format
			}
			if opt.formatFile != nil {
				formatFile = *opt.formatFile
			}
			if opt.color != nil {
				color = *opt.color
			}
			if opt.backColor != nil {
				backColor = *opt.backColor
			}
			if opt.file != nil {
				l.file = opt.file
			}
		}
	}

	l.format = initForm(format, levelKey, color, backColor)
	l.formatFile = initForm(formatFile, levelKey, color, backColor)
	return l
}

//initForm парсит основное форматирование уровня, заполняя управляющими байтами
func initForm(format string, level levelKey, color, bColor int) string {
	colorF := "\033[38;5;" + fmt.Sprint(color) + "m"
	colorB := "\033[48;5;" + fmt.Sprint(bColor) + "m"
	format = strings.ReplaceAll(format, "{{message}}", "%v")
	format = strings.ReplaceAll(format, "{{level}}", string(level))
	format = strings.ReplaceAll(format, "{{color:f}}", colorF)
	format = strings.ReplaceAll(format, "{{color:b}}", colorB)
	format = strings.ReplaceAll(format, "{{color:rf}}", "\033[39m")
	format = strings.ReplaceAll(format, "{{color:rb}}", "\033[49m")
	format = strings.ReplaceAll(format, "{{r}}", "\033[0m")
	format = strings.ReplaceAll(format, "{{font:b}}", "\033[1m")
	format = strings.ReplaceAll(format, "{{font:rb}}", "\033[22m")
	return format
}

//parser парсит форматирование сообщения, заполняя необходимыми данными
func (lv *level) parser(format string) string {
	regex := regexp.MustCompile("{{time[^{}]+")
	dateFind := regex.FindString(format)
	regex = regexp.MustCompile("{{file[^{}]+")
	fileFind := regex.FindString(format)
	dateForm := strings.TrimPrefix(dateFind, "{{time:")
	fileCalldeph, _ := strconv.Atoi(strings.TrimPrefix(fileFind, "{{file:"))
	t := time.Now()
	tForm := t.Format("01 02 2006:04:05.000000")
	if dateForm != "" {
		tForm = t.Format(dateForm)
	}
	if dateFind == "" {
		dateFind = "{{time"
	}
	if fileFind == "" {
		fileFind = "{{file"
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
	format = strings.ReplaceAll(format, dateFind+"}}", tForm)
	format = strings.ReplaceAll(format, fileFind+"}}", strings.Join(trace, "\n\t⮤ "))
	return format
}
