// Copyright (c) 2017 Kuguar <licenses@kuguar.io> Author: Adrian P.K. <apk@kuguar.io>
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	// UNSPECIFIED logs nothing
	UNSPECIFIED Level = iota // 0 :
	// DEBUG logs everything
	DEBUG // 1
	// TRACE logs everything
	TRACE // 2
	// INFO logs Info, Warnings and Errors
	INFO // 3
	// WARNING logs Warning and Errors
	WARNING // 4
	// ERROR just logs Errors
	ERROR // 5
)

// Level holds the log level.
type Level int

// Package level variables, which are pointer to log.Logger.
var (
	DebugLogger   *log.Logger
	TraceLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

// initLog initializes log.Logger objects
func initLog(
	debugHandle io.Writer,
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	isFlag bool) {

	// Flags for defines the logging properties, to log.New
	flag := 0
	if isFlag {
		flag = log.Ldate | log.Ltime | log.Lmicroseconds //| log.Lshortfile
	}

	// Create log.Logger objects.
	DebugLogger = log.New(debugHandle, "DEBUG: ", flag)
	TraceLogger = log.New(traceHandle, "TRACE: ", flag)
	InfoLogger = log.New(infoHandle, "INFO: ", flag)
	WarningLogger = log.New(warningHandle, "WARNING: ", flag)
	ErrorLogger = log.New(errorHandle, "ERROR: ", flag)
}

// SetLogLevel sets the logging level preference
func SetLogLevel(fullLogPath string, level Level) {
	//fullLogPath := path.Join(base, "logs/logs.txt")
	// Creates os.*File, which has implemented io.Writer intreface
	f, err := os.OpenFile(fullLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err.Error())
	}
	// Calls function initLog by specifying log level preference.
	switch level {
	case DEBUG:
		initLog(f, f, f, f, f, true)
		return

	case TRACE:
		initLog(ioutil.Discard, f, f, f, f, true)
		return

	case INFO:
		initLog(ioutil.Discard, ioutil.Discard, f, f, f, true)
		return

	case WARNING:
		initLog(ioutil.Discard, ioutil.Discard, ioutil.Discard, f, f, true)
		return

	case ERROR:
		initLog(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard, f, true)
		return

	default:
		initLog(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard, false)
		f.Close()
		return
	}
}

// Debug - calls l.Output to print to the logger.
func Debug(logText string) {
	DebugLogger.Println(logText)
}

// Debugf - calls l.Output to print to the logger.
func Debugf(format string, v ...interface{}) {
	DebugLogger.Printf(format, v)
}

// Trace - calls l.Output to print to the logger.
func Trace(logText string) {
	TraceLogger.Println(logText)
}

// Tracef - calls l.Output to print to the logger.
func Tracef(format string, v ...interface{}) {
	TraceLogger.Printf(format, v)
}

// Info - calls l.Output to print to the logger.
func Info(logText string) {
	InfoLogger.Println(logText)
}

// Infof - calls l.Output to print to the logger.
func Infof(format string, v ...interface{}) {
	InfoLogger.Printf(format, v)
}

// Warn - calls l.Output to print to the logger.
func Warn(logText string) {
	WarningLogger.Println(logText)
}

// Warnf - calls l.Output to print to the logger.
func Warnf(format string, v ...interface{}) {
	WarningLogger.Printf(format, v)
}

// Error - calls l.Output to print to the logger.
func Error(logText string) {
	ErrorLogger.Println(logText)
}

// Errorf - calls l.Output to print to the logger.
func Errorf(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v)
}

// Dump - calls l.Output to print error to the logger.
func Dump(error error) {
	ErrorLogger.Println(error.Error())
}

// Fatal - calls l.Output to print error to the logger and call os.Exit(1).
func Fatal(error error) {
	ErrorLogger.Fatal(error.Error())
}
