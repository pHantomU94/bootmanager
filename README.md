# **Bootmanager（脚本启动管理器)**

Bootmanager 是一个支持简单并行批量执行脚本的脚本启动器

## **目录**

### [三种模式](#三种模式)
#### [预设模式](#预设模式)
#### [自定义命令模式](#自定义命令模式)
#### [自定义格式模式](#自定义格式模式)
### [配置文件说明](#配置文件说明)
### [执行脚本存放路径说明](#执行脚本存放路径说明)
### [全部参数说明](#全部参数说明)

### **三种模式**
#### **预设模式**

```bash
# 执行配置文件中的boot预设
bootmanager -b 
# 执行配置文件中的config预设
bootmanager -c 
# 执行配置文件中的send预设
bootmanager -s 
```

#### **自定义命令模式**
```bash
# 执行配置文件中的[option]命令
bootmanager -o [option]
```
如，配置文件如下
```
{
    "boot": {
        "pattern": "boot.py",
        "interpreter": "python3",
        "parallel": true
    },
    "config": {
        "pattern": "config_board.py",
        "interpreter": "python3",
        "parallel": true
    },
    "send": {
        "pattern": "send_valid_data.py",
        "interpreter": "python3",
        "parallel": true
    },
    "custom": {
        "pattern": "custom.py",
        "interpreter": "python3",
        "parallel": true
    },
    "test": {
        "pattern": "test.py",
        "interpreter": "python3",
        "parallel": true
    }
}
```
```bash
# 批量执行test.py脚本
bootmanager -o test 
# 批量执行custom.py脚本
bootmanager -o custom 
```

#### **自定义格式模式**

```bash
# 执行对应pattern的脚本，使用解析器为interpreter
bootmanager -p [pattern] -i [interpreter]
```
```bash
# e.g. 可以直接并行执行python test*.py
bootmanager -p test.py -i python 
```

>NOTE:该模式下只支持并行执行，且必须指定interpreter，且自定义的pattern执行时优先级高于同时添加的前两种模式的操作

### **配置文件说明**

json格式，每个键值对采用以下形式
```
[option]: {
	"pattern": [scriptsPattern],
	"interpreter": [interpreter],
	"parallel": [bool]
}
option，string类型，表示操作的名称
pattern，string类型，表示脚本的格式
interpreter，string类型，表示脚本的解析器
parallel，bool类型值，表示是否并行
```

### **执行脚本存放路径说明**

如果使用`-d`参数指定执行的脚本路径，则需要将`pattern*.suffix`样式的脚本存放在该路径下，如果不加`-d`则指代当前目录

`pattern*.suffix`中`*`代表执行的数字`（0~∞）`

### **全部参数说明**
```bash
A convenient parallel scripts executor built with
          love by yinshijun in Go.
          You can use it to easily execute a series of parallel scripts.

Usage:
  bootmanager [flags]

Flags:
  -b, --boot                 Only boot the specified boards
  -f, --config-file string   Config file (default is /usr/local/bootmanager/config.json) (default "/usr/local/bootmanager/config.json")
  -c, --configure            Only configure the specified boards
  -h, --help                 help for bootmanager
  -i, --interpreter string   Specify interpreter for the custom option
  -l, --log                  Save log file
  -n, --numbers string       Specify the scripts numbers
  -o, --option string        Specify custom option
  -p, --pattern string       Specify custom script pattern
  -s, --send                 Use Viper for Only Start the sending data program of the specified board
  -d, --workdir string       Work directory (default is current directory)
```