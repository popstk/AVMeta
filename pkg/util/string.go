package util

import (
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// GetCode 从文件中提取番号信息
//
// filename 字符串，传入要提取的文件名称，
// filter 字符串，要对文件名称进行过滤的规则信息。
func GetCode(filename string, regs, filters []string) string {
	// 获取正确文件名
	filename = filepath.Base(strings.ToLower(filename))
	// 删除扩展名
	filename = strings.TrimSuffix(filename, path.Ext(filename))
	// 转为小写
	filename = strings.ToLower(filename)

	// 优先提取
	for _, reg := range regs {
		val := regexp.MustCompile(reg).FindString(filename)
		if len(val) > 0 {
			return val
		}
	}

	// 循环过滤
	for _, f := range filters {
		// 过滤
		filename = strings.ReplaceAll(filename, f, "")
	}
	// 将所有 . 替换为 -
	filename = strings.ReplaceAll(filename, ".", "-")
	// 过滤空格
	filename = strings.TrimSpace(filename)

	return filename
}

// GetNumberPath 通过配置信息，获取到正确的保存路径
//
// replaceStr map对象，通过转换后的媒体各项数据，
// cfg 配置信息，用以读取保存路径规则。
func GetNumberPath(replaceStr map[string]string, cfg *ConfigStruct) string {
	// 获取保存规则
	rule := cfg.Path.Directory
	// 循环替换
	for key, val := range replaceStr {
		rule = strings.ReplaceAll(rule, key, val)
	}

	// 定义特殊字符数组
	filter := []string{"\\", ":", "*", "?", `"`, "<", ">", "|"}
	// 循环过滤
	for _, v := range filter {
		rule = strings.ReplaceAll(rule, v, "")
	}
	// 多余的反斜线
	rule = strings.ReplaceAll(rule, "//", "/")

	return filepath.Join(GetRunPath(), cfg.Path.Success, rule)
}

// CheckDomainPrefix 检查域名最后是否存在斜线并返回无斜线域名
//
// domain 字符串，传入域名
func CheckDomainPrefix(domain string) string {
	// 是否为空
	if domain == "" {
		return ""
	}

	// 获取最后一个字符
	last := domain[len(domain)-1:]
	// 如果是斜线
	if last == "/" {
		domain = domain[:len(domain)-1]
	}

	return domain
}

// IntroFilter 简介信息过滤
func IntroFilter(intro string) string {
	// 替换<br>
	intro = strings.ReplaceAll(intro, "<br>", "\n")
	intro = strings.ReplaceAll(intro, "<br/>", "\n")
	intro = strings.ReplaceAll(intro, "<br />", "\n")
	// 替换\r\n
	intro = strings.ReplaceAll(intro, "\r\n", "\n")
	// 替换\r
	intro = strings.ReplaceAll(intro, "\r", "\n")
	// 替换\n\n
	intro = strings.ReplaceAll(intro, "\n\n", "\n")

	// 清除多余空白
	return strings.TrimSpace(intro)
}
