package logger

import (
	"avitoTelegram/config"
	"fmt"
	"github.com/gookit/color"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func LogRequest(r *http.Request) {
	if !config.Logging {
		return
	}
	params := getRequestParams(r.URL.Query())

	currentTime := time.Now().Format("15:04:05")
	fmt.Printf("\n******************** request (%s) ********************\n", currentTime)
	fmt.Println("method:", r.Method)
	fmt.Println("address:", r.URL.Path)
	if len(params) != 0 {
		fmt.Println("params:")
		for k, v := range params {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
	fmt.Println("************************************************************")
}

func LogResponse(statusCode int) {
	if !config.Logging {
		return
	}
	currentTime := time.Now().Format("15:04:05")
	pc, file, line, ok := runtime.Caller(2)

	printColor := color.Style{}
	if statusCode < 200 {
		printColor = color.Style{color.Cyan}
	} else if statusCode >= 200 && statusCode < 300 {
		printColor = color.Style{color.Green}
	} else if statusCode >= 300 && statusCode < 400 {
		printColor = color.Style{color.Magenta}
	} else if statusCode >= 400 && statusCode < 500 {
		printColor = color.Style{color.Yellow}
	} else if statusCode >= 500 {
		printColor = color.Style{color.FgRed}
	}
	printColor.Printf("\n******************* response (%s) ********************\n", currentTime)
	if ok {
		fileName := strings.Split(filepath.Base(file), ".")[0]
		funcName := strings.Replace(runtime.FuncForPC(pc).Name(), ".", "/"+fileName+".", -1)
		funcName = substringAfter(funcName, "/")
		fmt.Printf("Response from function - %s, line %d\n", funcName, line)
	}
	printColor.Println("status code: ", statusCode)
	printColor.Println("************************************************************")
}

func LogResponseError(err error, statusCode int) {
	if !config.Logging {
		return
	}
	pc, file, line, ok := runtime.Caller(2)
	currentTime := time.Now().Format("15:04:05")
	color.Red.Printf("\n***************** response error (%s) ****************\n", currentTime)
	if ok {
		fileName := strings.Split(filepath.Base(file), ".")[0]
		funcName := strings.Replace(runtime.FuncForPC(pc).Name(), ".", "/"+fileName+".", -1)
		funcName = substringAfter(funcName, "/")
		fmt.Printf("In function - %s, on line %d error:\n", funcName, line)
	}
	color.Red.Println(err)
	color.Red.Println("status code: ", statusCode)
	color.Red.Println("************************************************************")
}

func Log(info interface{}) {
	if !config.Logging {
		return
	}
	pc, file, line, ok := runtime.Caller(1)
	currentTime := time.Now().Format("15:04:05")

	color.Cyan.Printf("\n******************** log info (%s) *******************\n", currentTime)

	if ok {
		fileName := strings.Split(filepath.Base(file), ".")[0]
		funcName := strings.Replace(runtime.FuncForPC(pc).Name(), ".", "/"+fileName+".", -1)
		funcName = substringAfter(funcName, "/")
		fmt.Printf("Called from function - %s, line %d:\n", funcName, line)
	}
	color.Cyan.Println(info)
	color.Cyan.Println("************************************************************")
}

func LogError(err error) {
	if !config.Logging {
		return
	}
	pc, file, line, ok := runtime.Caller(1)
	currentTime := time.Now().Format("15:04:05")
	color.Red.Printf("\n********************* error (%s) *********************\n", currentTime)
	if ok {
		fileName := strings.Split(filepath.Base(file), ".")[0]
		funcName := strings.Replace(runtime.FuncForPC(pc).Name(), ".", "/"+fileName+".", -1)
		funcName = substringAfter(funcName, "/")
		fmt.Println("In function -", funcName, ", on line", line, "error:")
	}
	color.Red.Println(err)
	color.Red.Println("************************************************************")
}

func LogErrorIf(err error) {
	if !config.Logging {
		return
	}
	if err == nil {
		return
	}
	pc, file, line, ok := runtime.Caller(1)
	currentTime := time.Now().Format("15:04:05")
	color.Red.Printf("\n********************* error (%s) *********************\n", currentTime)
	if ok {
		fileName := strings.Split(filepath.Base(file), ".")[0]
		funcName := strings.Replace(runtime.FuncForPC(pc).Name(), ".", "/"+fileName+".", -1)
		funcName = substringAfter(funcName, "/")
		fmt.Println("In function -", funcName, ", on line", line, "error:")
	}
	color.Red.Println(err)
	color.Red.Println("************************************************************")
}

func LogFatal(err error) {
	if !config.Logging {
		return
	}
	pc, file, line, ok := runtime.Caller(1)
	currentTime := time.Now().Format("15:04:05")
	color.Red.Printf("\n****************** fatal error (%s) ******************\n", currentTime)
	if ok {
		fileName := strings.Split(filepath.Base(file), ".")[0]
		funcName := strings.Replace(runtime.FuncForPC(pc).Name(), ".", "/"+fileName+".", -1)
		funcName = substringAfter(funcName, "/")
		fmt.Println("In function -", funcName, ", on line", line, "fatal error:")
	}
	color.Red.Println(err)
	color.Red.Println("************************************************************")
	os.Exit(2)
}

func LogFatalIf(err interface{}) {
	if !config.Logging {
		return
	}
	if err == nil {
		return
	}
	pc, file, line, ok := runtime.Caller(1)
	currentTime := time.Now().Format("15:04:05")
	color.Red.Printf("\n****************** fatal error (%s) ******************\n", currentTime)
	if ok {
		fileName := strings.Split(filepath.Base(file), ".")[0]
		funcName := strings.Replace(runtime.FuncForPC(pc).Name(), ".", "/"+fileName+".", -1)
		funcName = substringAfter(funcName, "/")
		fmt.Printf("In function - %s, on line %d error:\n", funcName, line)
	}
	color.Red.Println(err)
	color.Red.Println("************************************************************")
	os.Exit(2)
}

func RollbackLog(err error) {
	if !config.Logging {
		return
	}
	pc, file, line, ok := runtime.Caller(1)
	currentTime := time.Now().Format("15:04:05")
	color.Red.Printf("************** transaction rollback (%s) *************\n", currentTime)
	if ok {
		fileName := strings.Split(filepath.Base(file), ".")[0]
		funcName := strings.Replace(runtime.FuncForPC(pc).Name(), ".", "/"+fileName+".", -1)
		funcName = substringAfter(funcName, "/")
		fmt.Printf("In transaction - %s, on line %d error:\n", funcName, line)
	}
	color.Red.Println(err)
	color.Red.Println("************************************************************")
}

func getRequestParams(values url.Values) map[string]string {
	urlValues := make(map[string]string, len(values))
	for k, v := range values {
		//fmt.Printf("key[%s] value[%s]\n", k, v[0])
		urlValues[k] = v[0]
	}
	//fmt.Println(urlValues)
	return urlValues
}

func substringAfter(str string, char string) string {
	return str[strings.IndexByte(str, char[0]):]
}
