package logger

import (
	"fmt"
	"os"
)

func InitOptions() Options {
	return &option{
		opt: map[levelKey]*levelOpt{
			LInfo:    {},
			LService: {},
			LWarning: {},
			LError:   {},
			LFatal:   {},
		},
	}
}

//SetLevelFormatConsole изменение формата вывода в консоль по конкретному уровню
func (o *option) SetLevelFormatConsole(level levelKey, format string) Options {
	if lv, ok := o.opt[level]; ok {
		lv.format = &format
	}
	return o
}

//SetLevelFormatFile изменение формата вывода в фаил по конкретному уровню
func (o *option) SetLevelFormatFile(level levelKey, format string) Options {
	if lv, ok := o.opt[level]; ok {
		lv.formatFile = &format
	}
	return o
}

//SetLevelFile изменение файла вывода по конкретному уровню
func (o *option) SetLevelFile(level levelKey, path string) Options {
	if lv, ok := o.opt[level]; ok {
		outFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			fmt.Printf("cannot crate file: %v\n", err)
		}
		lv.file = outFile
	}
	return o
}

//SetLevelColor изменение цвета сообщения по конкретному уровню
func (o *option) SetLevelColor(level levelKey, color int) Options {
	if lv, ok := o.opt[level]; ok {
		lv.color = &color
	}
	return o
}

//SetLevelBackColor изменение цвета фона сообщения по конкретному уровню
func (o *option) SetLevelBackColor(level levelKey, color int) Options {
	if lv, ok := o.opt[level]; ok {
		lv.backColor = &color
	}
	return o
}

//SetFormatConsole изменение формата вывода в консоль для всех уровней
func (o *option) SetFormatConsole(format string) Options {
	for _, v := range o.opt {
		v.format = &format
	}
	return o
}

//SetFormatFile изменение формата вывода в фаил  для всех уровней
func (o *option) SetFormatFile(format string) Options {
	for _, v := range o.opt {
		v.formatFile = &format
	}
	return o
}

//SetFile изменение файла вывода  для всех уровней
func (o *option) SetFile(path string) Options {
	outFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		fmt.Printf("cannot crate file: %v\n", err)
	}
	for _, v := range o.opt {
		v.file = outFile
	}
	return o
}

//Get
func (o *option) get(key levelKey) *levelOpt {
	if level, ok := o.opt[key]; ok {
		return level
	}
	return nil
}
