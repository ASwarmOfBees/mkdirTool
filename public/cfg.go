package public

import (
	"container/list"
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"sync"

	"mkdirTool/tool"
	"mkdirTool/util"
)

type GlobalConfig struct {
	IgoreFile []string `json:"ignoreFile"`
	IgoreDir  []string `json:"igoreDir"`
	SaveFile  string   `json:"resultFile"`
}

var (
	ConfigFile string
	Config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

var WaitGroup util.WaitGroupWrapper //异步线程
var DirQueue Queue                  //目录集合

func Init() {
	DirQueue = Queue{dirs: list.New()}
}

//解析配置文件cfg.json
func ParseConfig(cfg string) {

	ConfigFile = cfg //保存配置文件路径

	configContent, openError := tool.ReadFile(cfg)
	if openError != nil {
		log.Fatalln("read config file:", cfg, "fail:", openError)
	}
	//	configStr := Bytes2String(configContent)
	//	configStr = strings.Trim(configStr, " ")

	var c GlobalConfig                               //定义了一个GlobalConfig结构 和配置文件中的json对应
	err := json.Unmarshal([]byte(configContent), &c) //解析json串
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}
	//	log.Println(Bytes2String(configContent))
	lock.Lock()         //加锁
	defer lock.Unlock() //超出作用域时 解锁  好处等同于域智能指针

	Config = &c //保存配置信息

	log.Println("read config file:", cfg, "successfully")
}

type State int

// iota 初始化后会自动递增
const (
	CaseSensitive   State = iota // value --> 0
	UnCaseSensitive              // value --> 1
	Wildcard
)

//是否忽略目录
func IgnoreDir(dir string, s State) bool {
	lock.Lock()
	defer lock.Unlock()

	var result bool
	if s == Wildcard {
		reg := regexp.MustCompile(strings.Join(Config.IgoreFile, "|"))
		result = reg.MatchString(dir)
	} else {
		for index := range Config.IgoreDir {
			switch s {
			case CaseSensitive:
				result = Config.IgoreDir[index] == dir
			case UnCaseSensitive:
				result = strings.EqualFold(Config.IgoreDir[index], dir)
				if result {
					return true
				}
			}
		}
	}

	return false
}

//是否忽略文件
func IgnoreFile(file string, s State) bool {
	lock.Lock()
	defer lock.Unlock()

	var result bool
	if s == Wildcard {
		reg := regexp.MustCompile(strings.Join(Config.IgoreFile, "|"))
		result = reg.MatchString(file)
	} else {
		for index := range Config.IgoreFile {
			switch s {
			case CaseSensitive:
				result = Config.IgoreFile[index] == file
			case UnCaseSensitive:
				result = strings.EqualFold(Config.IgoreFile[index], file)
				if result {
					return true
				}
			}
		}
	}

	return false
}
