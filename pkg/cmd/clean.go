package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (e *Executor) initClean() {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "文件和目录情况",
		Long:  `按照规则清理无用的目录`,
		Run:   e.clean,
	}
	cmd.Flags().BoolP("force", "f", false, "执行清理")
	e.rootCmd.AddCommand(cmd)
}

func (e *Executor) clean(c *cobra.Command, _ []string) {
	e.initLog()
	e.initConfig()

	force, err := c.Flags().GetBool("force")
	if err != nil {
		log.Fatal(err)
	}

	dir := e.WorkPath()
	log.Infof("扫描目录 %s", dir)

	// find non-empty dir
	nonEmptyDir := make(map[string]struct{})

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
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

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		// 小于1m且不是特定扩展名就可以删除
		if fileInfo.Size() < 1024*1024 {
			ext := strings.ToLower(filepath.Ext(d.Name()))
			if _, ok := videoExts[ext]; !ok {
				return nil
			}
		}

		relDir := filepath.Dir(rel)
		for _, ok := nonEmptyDir[relDir]; !ok && relDir != "."; {
			nonEmptyDir[relDir] = struct{}{}
			relDir = filepath.Dir(relDir)
		}

		return nil
	})

	if err != nil {
		log.Errorf("walk dir err:%v", err)
		return
	}

	log.Debugf("nonEmptyDir: %+v", nonEmptyDir)

	var removeDir []string
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir == path {
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if e.InIgnoreDir(rel) {
			return filepath.SkipDir
		}

		log.Debugf("rel %s", rel)

		_, ok := nonEmptyDir[rel]
		if ok {
			return nil
		}

		removeDir = append(removeDir, path)
		return nil
	})

	if err != nil {
		log.Errorf("walk dir err:%v", err)
		return
	}

	if len(removeDir) == 0 {
		return
	}

	for _, p := range SortDir(removeDir) {
		log.Infof("remove dir %s", p)
		if force {
			if err := os.RemoveAll(p); err != nil {
				log.Errorf("err: %v", err)
			}
		}
	}
}

func SortDir(list []string) []string {
	if len(list) <= 1 {
		return list
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})

	ret := []string{list[0]}
	last := list[0]
	for i := 1; i < len(list); i++ {
		val := list[i]
		if strings.HasPrefix(val, last) {
			continue
		}

		last = val
		ret = append(ret, val)
	}

	return ret
}
