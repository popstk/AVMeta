package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ylqjgm/AVMeta/pkg/media"
	"github.com/ylqjgm/AVMeta/pkg/util"
	"io/fs"
	"path/filepath"
	"strings"
)

func (e *Executor) initRoot() {
	cmd := &cobra.Command{
		Use:   "AVMeta",
		Short: "一款使用 Golang 编写的跨平台 AV 元数据刮削器",
		Long: `AVMeta 是一款使用 Golang 编写的跨平台 AV 元数据刮削器
使用 AVMeta, 您可自动将 AV 电影进行归类整理
并生成对应媒体库元数据文件`,
		Run: e.rootRunFunc,
	}

	cmd.PersistentFlags().StringVarP(&e.workPath, "path", "p", "", "设置扫描目录")
	cmd.PersistentFlags().StringVarP(&e.configFile, "config", "c", "", "指定配置文件")
	cmd.PersistentFlags().BoolVarP(&e.verbose, "verbose", "v", false, "详细模式")
	e.rootCmd = cmd
}

// root命令执行函数
func (e *Executor) rootRunFunc(_ *cobra.Command, _ []string) {
	e.initLog()
	e.initConfig()

	files := e.findFiles()
	log.Infof("共探索到 %d 个视频文件, 开始刮削整理...", len(files))

	for i, file := range files {
		e.packProcess(i, file)
	}
}

func (e *Executor) findFiles() []string {
	dir := e.WorkPath()
	log.Infof("扫描目录 %s", dir)

	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			if e.InIgnoreDir(rel) {
				return filepath.SkipDir
			}

			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if _, ok := videoExts[ext]; ok {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return files
}

// 刮削进程
func (e *Executor) packProcess(i int, file string) {
	// 刮削整理
	m, err := media.Pack(file, e.cfg)
	// 检查
	if err != nil {
		// 恢复文件
		if len(e.cfg.Path.Fail) > 0 {
			util.FailFile(file, e.cfg.Path.Fail)
		}

		log.Errorf("pack err: %v", err)
		return
	}

	// 输出正确
	log.Infof("[%d]Done[%s] -> [%s]", i+1, m.Source, m.DirPath)
}
