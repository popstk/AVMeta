package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ylqjgm/AVMeta/pkg/util"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (e *Executor) initClean() {
	e.rootCmd.AddCommand(&cobra.Command{
		Use:   "clean",
		Short: "文件和目录情况",
		Long:  `按照规则清理无用的目录`,
		Run:   e.clean,
	})
}

func (e *Executor) clean(cmd *cobra.Command, args []string) {
	log.SetLevel(log.DebugLevel)

	dir, err := cmd.Flags().GetString("p")
	if err != nil {
		log.Errorf("get p err: %v", err)
		return
	}

	if len(dir) == 0 {
		dir = util.GetRunPath()
	}

	log.Infof("clean dir %s", dir)

	// find non-empty dir
	nonEmptyDir := make(map[string]struct{})

	ignoreDir := map[string]struct{}{
		"success": {},
		"fail":    {},
	}

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			if _, ok := ignoreDir[rel]; ok {
				return filepath.SkipDir
			}

			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if _, ok := videoExts[ext]; !ok {
			return nil
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

		if _, ok := ignoreDir[rel]; ok {
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

	sort.Slice(removeDir, func(i, j int) bool {
		return len(removeDir[i]) > len(removeDir[j])
	})

	for _, p := range removeDir {
		log.Infof("remove dir %s", p)
		if err := os.RemoveAll(p); err != nil {
			log.Errorf("err: %v", err)
		}
	}
}
