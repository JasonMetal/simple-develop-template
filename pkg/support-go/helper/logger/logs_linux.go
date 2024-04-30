//go:build !windows

package logger

import (
	"os"
	"strings"
	"syscall"
)

// 崩溃等标准输出日志重定向到文件
func RedirectLog() {
	pro := ""
	path, err := os.Getwd()
	if err == nil {
		paths := strings.Split(path, "/")

		pro = paths[len(paths)-1]
	}

	logFile, _ := os.OpenFile("/data/log/"+pro, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	//syscall.Dup2(int(logFile.Fd()), 1)
	syscall.Dup2(int(logFile.Fd()), 2)
}
