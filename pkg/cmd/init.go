package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"os"
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
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		e.workPath = dir
	}

	return e.workPath
}

func (e *Executor) initLog() {
	if e.verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
		jww.SetStdoutThreshold(jww.LevelTrace)
	}
}
