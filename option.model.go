package logger

import "os"

type Options interface {
	SetLevelFormatConsole(level levelKey, format string) Options
	SetLevelFormatFile(level levelKey, format string) Options
	SetLevelFile(level levelKey, path string) Options
	SetLevelColor(level levelKey, color int) Options
	SetLevelBackColor(level levelKey, color int) Options
	SetFormatConsole(format string) Options
	SetFormatFile(format string) Options
	SetFile(path string) Options

	get(key levelKey) *levelOpt
}

type option struct {
	opt map[levelKey]*levelOpt
}

//levelOpt модель поций каждого уровня логгера
type levelOpt struct {
	color      *int
	backColor  *int
	format     *string
	formatFile *string
	file       *os.File
}
