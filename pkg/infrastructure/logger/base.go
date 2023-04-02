package logger

import (
	"fmt"
	"location_store/pkg/utils"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SetupLogger ...
// Sets up the logger with preffered configurations according to the environment.
func SetupLogger() {
	// Global Time Format
	zerolog.TimeFieldFormat = "2 January 2006 Monday, 15:04:05 -07"

	// Global Path Format
	// Trims the file path, writes the path after backend folder and line.
	// If you want full path, just remove it. Default attitude is full path + line number
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		funcName := strings.Split(runtime.FuncForPC(pc).Name(), ".")[len(strings.Split(runtime.FuncForPC(pc).Name(), "."))-1]
		filePath := strings.Split(path.Dir(file), utils.ProjectFolderName)[1]
		fileName := path.Base(file)
		lineNumber := strconv.Itoa(line)
		return path.Join(funcName+" | "+filePath, fileName+":"+lineNumber)
	}

	// Deciding format of log output according to environment variable.
	// If the code runs on "LOCAL", writing the log output to console is pretty.
	// Else, output is in json format.
	if strings.Contains(utils.WorkEnvironment, "LOCAL") {
		log.Logger = prettyLoggerWithColor()
		// log.Logger = prettyLogger() // If your console is having color issues, try this.
	} else {
		log.Logger = prettyLogger()

		// Global Log Level
		// Least level of logs that is going to be printed out. "Trace" is the default value.
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

}

// Formatting and Customization for json output: Caller (address of log) and Timestamp added.
func jsonLogger() (jsonLogger zerolog.Logger) {
	jsonLogger = zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	return jsonLogger
}

// Formatting and Customization for console output.
// --- Colorless ---
func prettyLogger() (prettyLogger zerolog.Logger) {
	// Caller (address of log) and Timestamp added.
	prettyLogger = zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()

	// Output Format
	prettyLogger = prettyLogger.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    true, // Color Configuration
			TimeFormat: "15:04",
			FormatLevel: func(i interface{}) string {
				if i == nil {
					i = " log"
				}
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))

			},
			FormatCaller: func(i interface{}) string {
				return fmt.Sprintf("%s |", i)
			},
			FormatMessage: func(i interface{}) string {
				if i == nil {
					return ""
				}
				return fmt.Sprintf("Message: '%s' |", i)

			},
			FormatErrFieldName: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("%s: ", i))
			},
			FormatErrFieldValue: func(i interface{}) string {
				return fmt.Sprintf("%s |", i)
			},
			FormatFieldValue: func(i interface{}) string {
				return fmt.Sprintf("'%s' |", i)
			},
		})

	return prettyLogger
}

// Formatting and Customization for console output.
// --- Colorfull ---
func prettyLoggerWithColor() (prettyLogger zerolog.Logger) {
	// Caller (address of log) and Timestamp added.
	prettyLogger = zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()

	// Output Format
	prettyLogger = prettyLogger.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false, // Color Configuration
		},
	)

	return prettyLogger
}
