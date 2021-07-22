package manager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 判断格式的合法性并返回相应脚本列表
func patternValid(pattern string, numbers []int) ([]string, bool) {
	workDir := viper.GetString("workDir")
	s := strings.Split(pattern, ".")
	if len(s) < 2 {
		logrus.Errorln("Invalid file pattern of boot")
		os.Exit(1)
	}
	scripts := make([]string, 0, len(numbers))

	// 未输入参数全部执行
	if numbers == nil {
		fileName := fmt.Sprintf("%s*.%s", s[0], s[1])
		path := filepath.Join(workDir, fileName)
		scripts, err := filepath.Glob(path)
		if err != nil {
			logrus.Errorln(os.Stderr, err.Error())
			os.Exit(1)
		}
		return scripts, true
	}

	for _, number := range numbers {
		fileName := fmt.Sprintf("%s%d.%s", s[0], number, s[1])
		path := filepath.Join(workDir, fileName)
		if !fileEixst(path) {
			return nil, false
		}
		scripts = append(scripts, path)
	}

	return scripts, true
}

// 判断操作合法性并返回相应脚本列表
func optionValid(option string, numbers []int) ([]string, bool) {
	pattern := viper.GetString(option + ".pattern")
	return patternValid(pattern, numbers)
}

func runOption(ctx context.Context, interpreter string, script string, index int, args []string) error {
	// logrus.Infof("%s start\n", script)
	cmdargs := make([]string, 0, 1 + len(args))
	cmdargs = append(cmdargs, script)
	cmdargs = append(cmdargs, args...)
	commandline := exec.CommandContext(ctx, interpreter, cmdargs...)
	// DEBUG: 这里仅作为调试功能
	if viper.GetBool("logFlag") {
		logPath := fmt.Sprintf("log%d.txt", index+1)
		logfile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		defer logfile.Close()
		if err != nil {
			logrus.Errorln(os.Stderr, err)
			return err
		}
		commandline.Stdout = logfile
		commandline.Stderr = logfile
	} else {
		commandline.Stdout = os.Stdout
		commandline.Stderr = os.Stderr
	}
	err := commandline.Start()
	if err != nil {
		logrus.Warnln(script, ":", err)
		return err
	}
	err = commandline.Wait()
	if err != nil {
		logrus.Errorln(script, ":", err)
		return err
	}
	// logrus.Tracef("%s end\n", script)
	return nil
}

// 重试函数
func retryOption(ctx context.Context, interpreter string, script string, number int, retries int, args []string) (err error) {
	// 重试指定次数脚本
	logrus.Warnf("Retry %s for %d times\n", script, retries)
	for i:=0; i<retries ; i++ {
		err = runOption(ctx, interpreter, script, number, args)
		if err == nil {
			logrus.Infof("Retry %s success at %d time\n", script, i+1)
			return
		}
	}
	logrus.Errorf("Retry %s failed\n", script)
	return
}

// 执行并行操作
func parallelRunOption(ctx context.Context, interpreter string, scripts []string, retries int, args []string) {
	var lock sync.Mutex
	wg := sync.WaitGroup{}
	failedArr := make([]int, 0)
	failedScripts := make([]string, 0)
	for _, script := range scripts {
		wg.Add(1)
		go func(script string) {
			defer wg.Done()
			number_string := regexp.MustCompile(`[0-9]+\.`).FindString(script)
			number, _ := strconv.Atoi(strings.Split(number_string, ".")[0])
			err := runOption(ctx, interpreter, script, number, args)
			if err != nil {
				if retries != 0 {
					// 失败重试
					err = retryOption(ctx, interpreter, script, number, retries, args)
				}
				// 重试失败
				if err != nil {
					lock.Lock()
					// WARN: 这里是从1开始定义的
					failedArr = append(failedArr, number)
					failedScripts = append(failedScripts, script)
					lock.Unlock()
				}
			}
		}(script)
		time.Sleep(time.Duration(10) * time.Microsecond)
	}
	wg.Wait()

	if len(failedArr) != 0 {
		sort.Sort(sort.IntSlice(failedArr))
		logrus.Infof("Option Done. Total: %d. Failed: %d\n", len(scripts), len(failedArr))
		logrus.Infof("Failed nodes: %v\n", failedArr)
	} else {
		logrus.Infof("Option Done. Total: %d. All success!\n", len(scripts))
	}
}

func serialRunOptin(ctx context.Context, interpreter string, scripts []string, args []string) {
	for index, script := range scripts {
		runOption(ctx, interpreter, script, index, args)
	}
}
