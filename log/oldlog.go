package log

import (
	"errors"
	"fmt"
	"go-study/log/tracer"
	"io"
	"os"

	colorable "github.com/mattn/go-colorable"
	opentrace "github.com/opentracing/opentracing-go"
	logging "github.com/whyrusleeping/go-logging"
)

func init() {
	SetupLogging()
}

var ansiGray = "\033[0;37m"
var ansiBlue = "\033[0;34m"

// LogFormats defines formats for logging (i.e. "color")
var LogFormats = map[string]string{
	"nocolor": "%{time:2006-01-02 15:04:05.000000} %{level} %{module} %{longfile}: %{message}",
	"color": ansiGray + "%{time:15:04:05.000} %{color}%{level:5.5s} " + ansiBlue +
		"%{module:10.10s}: %{color:reset}%{message} " + ansiGray + "%{longfile}%{color:reset}",
	"file": "%{time:2006-01-02 15:04:05.000000} %{level} %{module} %{shortfile}: %{message}",
}

var defaultLogFormat = "color"
var defaultLogFileFormat = "file"

// Logging environment variables
const (
	envLogging     = "LW_LOGGING"
	envLoggingFmt  = "LW_LOGGING_FMT"
	envLoggingFile = "LW_LOGGING_FILE"
)

// ErrNoSuchLogger is returned when the util pkg is asked for a non existant logger
var ErrNoSuchLogger = errors.New("Error: No such logger")

// loggers is the set of loggers in the system
var loggers = map[string]*logging.Logger{}

// SetupLogging will initialize the logger backend and set the flags.
func SetupLogging() {

	lfmt := LogFormats[os.Getenv(envLoggingFmt)]
	if lfmt == "" {
		lfmt = LogFormats[defaultLogFormat]
	}

	backend := logging.NewLogBackend(colorable.NewColorableStderr(), "", 0)
	logging.SetBackend(backend)
	logging.SetFormatter(logging.MustStringFormatter(lfmt))

	lvl := logging.ERROR

	if logenv := os.Getenv(envLogging); logenv != "" {
		var err error
		lvl, err = logging.LogLevel(logenv)
		if err != nil {
			fmt.Println("error setting log levels", err)
		}
	}

	// TracerPlugins are instantiated after this, so use loggable tracer
	// by default, if a TracerPlugin is added it will override this
	lgblRecorder := tracer.NewLoggableRecorder()
	lgblTracer := tracer.New(lgblRecorder)
	opentrace.SetGlobalTracer(lgblTracer)

	// fmt.Println("setupLogging..............")
	SetAllLoggers(lvl)
}

// 设置日志文件
//   - 可以切割文件，根据文件内容大小和文件日期来进行切割
//     fileName 需要切割的日志文件名称
//     dateSlice 是否通过日期来进行分割 有效值分别为 y m d h
//     maxSize 是否通过文件大小进行分割，单位kb
//     enableConsole 是否允许控制台输出
func SetupLogFile(fileName, dateSlice string, maxSize int64, enableConsole bool, out ...io.Writer) {
	// f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	// 文件打开失败，直接使用之前默认控制台输出
	// 	return
	// }
	f := NewFileWrite(fileName, dateSlice, maxSize)
	backendFile := logging.NewLogBackend(f, "", 0)
	if enableConsole {
		var backendConsole *logging.LogBackend
		if len(out) > 0 {
			backendConsole = logging.NewLogBackend(out[0], "", 0)
		} else {
			backendConsole = logging.NewLogBackend(colorable.NewColorableStderr(), "", 0)
		}
		logging.SetBackend(backendFile, backendConsole)
	} else {
		logging.SetBackend(backendFile)
	}

	lfmt := LogFormats[os.Getenv(envLoggingFile)]
	if lfmt == "" {
		lfmt = LogFormats[defaultLogFileFormat]
	}
	logging.SetFormatter(logging.MustStringFormatter(lfmt))
}

// SetDebugLogging calls SetAllLoggers with logging.DEBUG
func SetDebugLogging() {
	SetAllLoggers(logging.DEBUG)
}

// SetAllLoggers changes the logging.Level of all loggers to lvl
func SetAllLoggers(lvl logging.Level) {
	logging.SetLevel(lvl, "")
	for n := range loggers {
		logging.SetLevel(lvl, n)
	}
}

// SetLogLevel changes the log level of a specific subsystem
// name=="*" changes all subsystems
func SetLogLevel(name, level string) error {
	lvl, err := logging.LogLevel(level)
	if err != nil {
		return err
	}

	// wildcard, change all
	if name == "*" {
		SetAllLoggers(lvl)
		return nil
	}

	// Check if we have a logger by that name
	if _, ok := loggers[name]; !ok {
		return ErrNoSuchLogger
	}

	logging.SetLevel(lvl, name)

	return nil
}

// GetSubsystems returns a slice containing the
// names of the current loggers
func GetSubsystems() []string {
	subs := make([]string, 0, len(loggers))

	for k := range loggers {
		subs = append(subs, k)
	}
	return subs
}

func getLogger(name string) *logging.Logger {
	log := logging.MustGetLogger(name)
	loggers[name] = log
	return log
}
