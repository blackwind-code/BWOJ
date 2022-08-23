package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/sirupsen/logrus"
)

func File2Str(fp string) string {
	tmp, e := ioutil.ReadFile(fp)
	if e != nil {
		log.Println("No such file!: "+fp+" Error msg: "+fp, e)
	}
	return strings.Trim(string(tmp), " \r\t\n")
}

// Str2File function
func Str2File(path, data string, permission uint32) error {
	return ioutil.WriteFile(path, []byte(data), os.FileMode(permission))
}

// Str2File function
func Str2FileAppend(path, data string, permission uint32) error {
	f, e := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.FileMode(permission))
	if e != nil {
		return e
	}
	defer f.Close()
	_, e = f.WriteString(data)
	if e != nil {
		return e
	}
	return nil
}

func GetCurrentFunctionInfo() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%s:%d", frame.File, frame.Function, frame.Line)
}

func Atof(value float64) string {
	return strconv.FormatFloat(value, 'f', 3, 64)
}

func MatchRegexPCRE2(regex string, target string) bool {
	re, e := regexp2.Compile(regex, 0)
	if e != nil {
		logrus.WithFields(
			logrus.Fields{
				"regex_regex":     regex,
				"regex_candidate": target,
				"error":           e,
				"where":           GetCurrentFunctionInfo(),
			}).Fatalf("Failed to compile Regex")
	}

	matched, e := re.MatchString(target)
	if e != nil {
		logrus.WithFields(
			logrus.Fields{
				"regex_regex":     regex,
				"regex_candidate": target,
				"error":           e,
				"where":           GetCurrentFunctionInfo(),
			}).Fatalf("Faled to match Regex")
	}

	return matched
}
