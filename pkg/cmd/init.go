package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ylqjgm/AVMeta/pkg/util"
)

func (e *Executor) initConfigFile() {
	e.rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "生成配置文件",
		Long:  `在当前目录下生成 config.yaml 配置文件`,
		Run: func(cmd *cobra.Command, args []string) {
			e.initConfig()
		},
	})
}

func (e *Executor) WorkPath() string {
	if len(e.workPath) == 0 {
		e.workPath = util.GetRunPath()
	}

	return e.workPath
}
