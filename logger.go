package main

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"local/global"
)

func setDefaultLogger() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = global.FormatTime("y-m-d h:i:s")
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: zerolog.TimeFieldFormat,
	})
}

// 根据配置文件设置logger
func setLogger() error {
	// 设置级别
	level := strings.ToLower(global.Config.Logger.Level)
	switch level {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "empty":
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
		return nil
	}

	// 设置时间格式
	if global.Config.Logger.TimeFormat == "timestamp" {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	} else {
		zerolog.TimeFieldFormat = global.FormatTime(global.Config.Logger.TimeFormat)
	}

	// 设置日志输出方式
	var output io.Writer
	var logFile *os.File
	var err error
	// 设置日志文件
	if global.Config.Logger.FilePath != "" {
		// 输出到文件
		if global.Config.Logger.FileMode == 0 {
			global.Config.Logger.FileMode = os.FileMode(0600)
		}
		logFile, err = os.OpenFile(global.Config.Logger.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, global.Config.Logger.FileMode)
		if nil != err {
			return err
		}
	}
	switch global.Config.Logger.Encode {
	// console编码
	case "console":
		if logFile != nil {
			output = zerolog.ConsoleWriter{
				Out:        logFile,
				NoColor:    true,
				TimeFormat: zerolog.TimeFieldFormat,
			}
		} else {
			output = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: zerolog.TimeFieldFormat,
			}
		}
	// json编码
	case "json":
		if logFile != nil {
			output = logFile
		} else {
			output = os.Stdout
		}
	default:
		return errors.New("从配置文件的logger.encode中获得了未知的参数，目前只支持json|console")
	}

	log.Logger = log.Output(output)

	return nil
}
