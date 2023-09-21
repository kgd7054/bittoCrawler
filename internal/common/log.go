package common

import (
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Config struct {
	UseTerminal        bool
	UseFile            bool
	TerminalJSONOutput bool
	VerbosityTerminal  int
	VerbosityFile      int
	FilePath           string
}

func InitLog(config Config) {
	hs := func() []log.Handler {
		h := []log.Handler{}
		if config.UseTerminal == true { //터미널 모드
			handler := log.LvlFilterHandler(log.Lvl(config.VerbosityTerminal),
				log.StreamHandler(os.Stdout, log.TerminalFormat(true)))
			h = append(h, handler)
		}
		if config.UseFile == true { //파일모드
			handler := log.LvlFilterHandler(log.Lvl(config.VerbosityFile), log.StreamHandler(&lumberjack.Logger{
				Filename:   config.FilePath,
				MaxSize:    64, // megabytes
				MaxBackups: 3,
				MaxAge:     28,   //days
				Compress:   true, // disabled by default
			}, log.JSONFormatOrderedEx(false, true)))
			h = append(h, handler)
		}
		return h
	}()
	log.Root().SetHandler(log.MultiHandler(hs...))
}

// Trace is a convenient alias for Root().Trace
func Trace(msg string, ctx ...interface{}) {
	log.Trace(msg, ctx...)
}

// Debug is a convenient alias for Root().Debug
func Debug(msg string, ctx ...interface{}) {
	log.Debug(msg, ctx...)
}

// Info is a convenient alias for Root().Info
func Info(msg string, ctx ...interface{}) {
	log.Info(msg, ctx...)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, ctx ...interface{}) {
	log.Warn(msg, ctx...)
}

// Error is a convenient alias for Root().Error
func Error(msg string, ctx ...interface{}) {
	log.Error(msg, ctx...)
}

// Crit is a convenient alias for Root().Crit
func Crit(msg string, ctx ...interface{}) {
	// util.SendChatAlert(msg)
	log.Crit(msg, ctx...)
}
