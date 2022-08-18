/*
Package cmd 命令行操作包。

AVMeta程序所有操作命令皆由此包定义，使用 cobra 第三方包编写。
*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ylqjgm/AVMeta/pkg/config"
	"runtime"
)

// 采集站点变量
var site string

// Executor 命令对象
type Executor struct {
	rootCmd *cobra.Command
	cfg     *config.Conf

	version   string
	commit    string
	built     string
	goVersion string
	platform  string

	workPath   string // 工作目录
	configFile string // 配置文件
	verbose    bool   // 详细模式

	ignoreDir map[string]struct{}
}

// NewExecutor 返回一个被初始化的命令对象。
//
// version 字符串参数，传入当前程序版本，
// commit 字符串参数，传入最后提交的 git commit，
// built 字符串参数，传入程序编译时间。
func NewExecutor(version, commit, built string) *Executor {
	e := &Executor{
		version:   version,
		commit:    commit,
		built:     built,
		goVersion: runtime.Version(),
		platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
	e.initRoot()
	e.initConfigFile()
	e.initActress()
	e.initNfo()
	e.initVersion()
	e.initClean()

	return e
}

// Execute 执行根命令。
func (e *Executor) Execute() error {
	return e.rootCmd.Execute()
}

// 初始化配置
func (e *Executor) initConfig() {
	// 获取配置
	cfg, err := config.GetConfig(e.configFile)
	// 检查
	if err != nil {
		log.Fatal(err)
	}

	// 配置信息
	e.cfg = cfg
	e.Init()
}

func (e *Executor) Init() {
	c := e.cfg
	ignoreDir := make(map[string]struct{})
	ignoreDir[c.Path.Success] = struct{}{}
	ignoreDir[c.Path.Fail] = struct{}{}
	e.ignoreDir = ignoreDir
}

func (e *Executor) InIgnoreDir(dir string) bool {
	if len(e.ignoreDir) > 0 {
		_, ok := e.ignoreDir[dir]
		return ok
	}

	return false
}
