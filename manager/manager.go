package manager

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 最大值
func max(a int ,b int ) int {
	if a > b {
		return a 
	} else {
		return b 
	}
}

// 去重合并
func merge(intervals [][2]int) []int{
	merged := make([][2]int, 0, len(intervals))

	//合并区间
	for _, interval := range intervals {
		if len(merged) == 0 || merged[len(merged)-1][1] < interval[0] {
			merged = append(merged, interval)
		} else {
			merged[len(merged)-1][1] = max(merged[len(merged)-1][1], interval[1])
		}
	}
	// 生成数组
	numbers := make([]int, 0, len(merged))
	for _, interval := range merged {
		for i:= interval[0]; i<=interval[1] ; i++ {
			numbers = append(numbers, i)
		}
	}
	return numbers
}

// 获取板子id
func getNum() ([]int, bool) {
	intervals := make([][2]int, 0, 0) 
	arg := viper.GetString("numbers")
	// 如果没有输入相关参数则直接返回
	if arg == "" {
		return nil, true
	}
	// 判断合法性
	valid, _ := regexp.MatchString(`(((\d+-\d+)|\d+),)*(((\d+-\d+)|\d+))?$`, arg)
	if !valid {
		return nil, false
	}
	match := regexp.MustCompile(`(((\d+-\d+)|\d+),)*(((\d+-\d+)|\d+))?$`).FindString(arg)
	if len(match) == 0 {
		return nil, false
	}
	// 提取子串
	subs := strings.Split(arg, ",")
	for _, sub := range subs {
		if strings.Contains(sub, "-") {
			nums := strings.Split(sub, "-")
			start, err := strconv.Atoi(nums[0])
			if err != nil {
				continue
			}
			end, err := strconv.Atoi(nums[1])
			if err != nil {
				continue
			}
			if start > end {
				continue
			}
			interval := [2]int{start, end}
			intervals = append(intervals, interval)

		} else {
			num, err := strconv.Atoi(sub)
			if err != nil {
				continue
			}
			interval := [2]int{num, num}
			intervals = append(intervals, interval)
		}
	}
	numbers := merge(intervals)
	// 数组去重
	return numbers, true
}

// 判断文件是否存在
func fileEixst(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// Run 执行程序入口
func Run(args []string) {
	numbers, valid := getNum()
	if !valid {
		fmt.Fprintln(os.Stderr, "Invalid parameter form for -n")
		os.Exit(1)
	}
	options := make([]string, 0, 3)

	// 生成操作表
	if viper.GetBool("bootFlag") {
		options = append(options, "boot")
	} 
	if viper.GetBool("configureFlag") {
		options = append(options, "config")
	}
	if viper.GetBool("sendFlag") {
		options = append(options, "send")
	}
	if viper.GetString("option") != "" {
		customOption := viper.GetString("option")
		if len(viper.GetStringMap(customOption)) == 0 {
			logrus.Warnf("No custom option [%s] in your configure file, please check.\n", customOption)
			return
		} 
		options = append(options, customOption)
	}

	customPattern := viper.GetString("pattern")
	customInterpreter := viper.GetString("interpreter")
	retries := viper.GetInt("retries")

	ctx := context.Background()
	if customPattern != "" {
		if customInterpreter == "" {
			logrus.Errorln("You need to specify the interpreter when using custom pattern.")
			os.Exit(1)
		}
		scripts, valid := patternValid(customPattern, numbers)
		if !valid {
			logrus.Errorln("Invalid board number range")
			os.Exit(1)
		}
		logrus.Infof("Custom pattern [%s] start.\n", customPattern)
		parallelRunOption(ctx, customInterpreter, scripts, retries, args)
	}

	// 判断是否有操作要执行
	if len(options) == 0 {
		// WARN: 这里需要按需修改为指定脚本类型或者从配置文件读取
		// options = append(options, "boot", "config", "send")
		// 这里无操作返回
		if customPattern == ""{
			logrus.Infoln("No option to do.")
		}
		return
	}

	// 脚本列表
	scripts_list := make([][]string, 0, len(options))
	// 判断操作是否合理
	for _, option := range options {
		scripts, valid := optionValid(option, numbers)
		if !valid {
			logrus.Errorln("Invalid board number range")
			os.Exit(1)
		}
		scripts_list = append(scripts_list, scripts)
	}
	// 按阶段执行操作
	for key, scripts := range scripts_list {
		option := options[key]
		interpreter := viper.GetString(option+".interpreter")
		parallel := viper.GetBool(option+".parallel")
		logrus.Infof("Option [%s] start.\n", option)
		if parallel {
			parallelRunOption(ctx, interpreter, scripts, retries, args)
		} else {
			serialRunOptin(ctx, interpreter, scripts, args)
		}
	}
}