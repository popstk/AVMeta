package cmd

import (
	"github.com/ylqjgm/AVMeta/pkg/logs"
	"path"

	"github.com/ylqjgm/AVMeta/pkg/media"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/spf13/cobra"
)

func (e *Executor) initRoot() {
	e.rootCmd = &cobra.Command{
		Use:   "AVMeta",
		Short: "一款使用 Golang 编写的跨平台 AV 元数据刮削器",
		Long: `AVMeta 是一款使用 Golang 编写的跨平台 AV 元数据刮削器
使用 AVMeta, 您可自动将 AV 电影进行归类整理
并生成对应媒体库元数据文件`,
		Run: e.rootRunFunc,
	}
	e.rootCmd.PersistentFlags().String("p", "", "设置目录")
}

// root命令执行函数
func (e *Executor) rootRunFunc(_ *cobra.Command, _ []string) {
	// 初始化日志
	logs.Log("logs")

	// 获取当前执行路径
	curDir := util.GetRunPath()
	logs.Info("walk dir %s", curDir)

	// 列当前目录
	files, err := util.WalkDir(curDir, e.cfg.Path.Success, e.cfg.Path.Fail)
	// 错误日志
	logs.FatalError(err)

	// 获取总量
	count := len(files)
	// 输出总量
	logs.Info("\n\n共探索到 %d 个视频文件, 开始刮削整理...\n", count)

	// 初始化进程
	wg := util.NewWaitGroup(2)

	// 循环视频文件列表
	for _, file := range files {
		// 计数加
		wg.AddDelta()
		// 刮削进程
		e.packProcess(file, wg)
	}

	// 等待结束
	wg.Wait()
}

// 刮削进程
func (e *Executor) packProcess(file string, wg *util.WaitGroup) {
	// 刮削整理
	m, err := media.Pack(file, e.cfg)
	// 检查
	if err != nil {
		// 恢复文件
		util.FailFile(file, e.cfg.Path.Fail)

		logs.Error("pack err: %v", err)
		// 进程
		wg.Done()

		return
	}

	// 输出正确
	logs.Info("文件 [%s] 刮削成功, 来源 [%s], 路径 [%s]", path.Base(file), m.Source, m.DirPath)

	// 进程
	wg.Done()
}
