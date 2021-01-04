package manager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 判断操作合法性并返回相应脚本列表
func optionValid(option string, numbers []int) ([]string, bool) {
	workDir := viper.GetString("workDir")
	patten := viper.GetString(option+".pattern")
	s := strings.Split(patten, ".")
	if len(s) < 2 {
		fmt.Fprintln(os.Stderr, "Invalid file pattern of boot")
		os.Exit(1)
	}
	scripts := make([]string, 0, len(numbers))

	// 未输入参数全部执行
	if numbers == nil {
		fileName := fmt.Sprintf("%s*.%s",s[0], s[1]) 
		path := filepath.Join(workDir, fileName)
		scripts, err := filepath.Glob(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		return scripts, true
	}

	for _, number := range numbers {
		fileName := fmt.Sprintf("%s%d.%s",s[0], number, s[1])
		path := filepath.Join(workDir, fileName)
		if !fileEixst(path) {
			return nil, false
		}
		scripts = append(scripts, path)
	}

	return scripts, true
}

func runOption(ctx context.Context, interpreter string, script string, index int) {
	logrus.Infof("%s start\n", script)
	commandline := exec.CommandContext(ctx, interpreter, script)

	// DEBUG: 这里仅作为调试功能
	if viper.GetBool("logFlag") {
		logPath := fmt.Sprintf("log%d.txt", index+1)
		logfile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		defer logfile.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		commandline.Stdout = logfile
		commandline.Stderr = logfile
	} else {
		commandline.Stdout = os.Stdout
		commandline.Stderr = os.Stderr
	}
	err := commandline.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	err = commandline.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	logrus.Infof("%s end\n", script)
}

// 执行并行操作
func parallelRunOption(ctx context.Context, interpreter string, scripts []string) {
	wg := sync.WaitGroup{}
	for index, script := range scripts {
		wg.Add(1)
		go func(script string) {
			defer wg.Done()
			runOption(ctx, interpreter, script, index)
		}(script)
		time.Sleep(time.Duration(10) * time.Microsecond)
	}
	wg.Wait()
}

func serialRunOptin(ctx context.Context, interpreter string, scripts []string) {
	for index, script := range scripts {
		runOption(ctx, interpreter, script, index)
	}
}