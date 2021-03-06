package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"bootmanager/manager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	wkDir string	// 工作令
	cfgFile string	// 配置文件
	number string	// 文件编号
	boot bool		// 仅启动
	confgure bool	// 仅配置
	send bool		// 仅发送
	log bool		// 保存日志
	pattern string  // 自定义格式
	option string 	// 自定义操作
	retries int		// 重传次数
	interpreter string
	rootCmd = &cobra.Command{
		Use: "bootmanager",
		Short: "Bootmanager is a parallel scripts boot entry",
		Long: `A convenient parallel scripts executor built with
	  love by yinshijun in Go.
	  You can use it to easily execute a series of parallel scripts.
	  Version: v 1.0.2`,
		// Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			manager.Run(args)
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	//解析参数
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config-file", "f", "/usr/local/bootmanager/config.json", "Config file (default is /usr/local/bootmanager/config.json)")
	rootCmd.PersistentFlags().BoolVarP(&boot, "boot", "b", false, "Only boot the specified boards")
	rootCmd.PersistentFlags().StringVarP(&wkDir, "workdir", "d", "", "Work directory (default is current directory)")
	rootCmd.PersistentFlags().StringVarP(&number, "numbers", "n", "", "Specify the scripts numbers")
	rootCmd.PersistentFlags().BoolVarP(&log, "log", "l", false, "Save log file")
	rootCmd.PersistentFlags().BoolVarP(&confgure, "configure", "c", false, "Only configure the specified boards")
	rootCmd.PersistentFlags().BoolVarP(&send, "send", "s", false, "Use Viper for Only Start the sending data program of the specified board")
	rootCmd.PersistentFlags().StringVarP(&pattern, "pattern", "p","", "Specify custom script pattern")
	rootCmd.PersistentFlags().StringVarP(&option, "option", "o","", "Specify custom option")
	rootCmd.PersistentFlags().StringVarP(&interpreter, "interpreter", "i","", "Specify interpreter for the custom option")
	rootCmd.PersistentFlags().IntVarP(&retries, "retries", "r", 0, "Specify retries for failed scripts")
	// 用viper收集参数
	viper.BindPFlag("bootFlag", rootCmd.PersistentFlags().Lookup("boot"))
	viper.BindPFlag("configureFlag", rootCmd.PersistentFlags().Lookup("configure"))
	viper.BindPFlag("sendFlag", rootCmd.PersistentFlags().Lookup("send"))
	viper.BindPFlag("workDir", rootCmd.PersistentFlags().Lookup("workdir"))
	viper.BindPFlag("numbers", rootCmd.PersistentFlags().Lookup("numbers"))
	viper.BindPFlag("pattern", rootCmd.PersistentFlags().Lookup("pattern"))
	viper.BindPFlag("option", rootCmd.PersistentFlags().Lookup("option"))
	viper.BindPFlag("interpreter", rootCmd.PersistentFlags().Lookup("interpreter"))
	viper.BindPFlag("logFlag", rootCmd.PersistentFlags().Lookup("log"))
	viper.BindPFlag("retries", rootCmd.PersistentFlags().Lookup("retries"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		_, err := os.Stat(cfgFile)
		if !(err == nil || os.IsExist(err)) {
			fmt.Fprintln(os.Stderr, "The config file is not exist.")
			os.Exit(1)
		}
	} else {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		viper.AddConfigPath(currentDir)
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		fileName := filepath.Join(currentDir, "config.json")
		_, err = os.Stat(fileName)
		if !(err == nil || os.IsExist(err)) {
			fmt.Fprintln(os.Stderr, "The config file is not exist. Please check /usr/local/bootmanager/config.json.")
			os.Exit(1)
		}
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	if wkDir == "" {
		wkDir, _ = os.Getwd()
	}
}

// Execute rootCmd执行入口
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// bootFile := viper.GetString("boot")
	// configFile := viper.GetString("config")
	// sendFile := viper.GetString("")
}