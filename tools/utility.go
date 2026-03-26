package tools

import (
	"fmt"
	"github.com/kpango/glg"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func GetFileLine() (string, int) {

	_, file, line, ok := runtime.Caller(2)
	if ok {
		return file, line
	}
	return "", -1
}

func IssueLog(message string, args ...interface{}) {

	var params []interface{}

	file, line := GetFileLine()
	params = append(params, goid())
	params = append(params, filepath.Base(file))
	params = append(params, line)

	for _, arg := range args {
		params = append(params, arg)
	}
	err := glg.Infof(message, params...)
	if err != nil {
		return
	}
	// log.Printf("[%d][%s]{%d} "+message, params...)
}

func goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func SimpleMergeMap(left, right map[string]string) *map[string]string {
	for key, rightVal := range right {
		left[key] = rightVal
	}
	return &left
}
