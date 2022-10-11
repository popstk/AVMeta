package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	UserAgent      = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Mobile Safari/537.36"
	DefaultTimeOut = 15 * time.Second
	RetryTime      = 3
)

// MakeRequest 创建一个远程请求对象
//
// 请求参数：
//
// method 字符串参数，传入请求类型，
// uri 字符串参数，传入请求地址，
// proxy 字符串参数，传入代理地址，
// body io读取接口，传入远程内容对象，
// header map结构，传入头部信息，
// cookies cookie数组，传入cookie信息。
//
// 返回数据：
//
// data 字节集，返回读取到的内容字节集，
// status 整数，返回请求状态码，
// err 错误信息。
func MakeRequest(method, uri, proxy string, body io.Reader, header map[string]string, cookies []*http.Cookie) (
	data []byte, status int, err error) {
	// 构建请求客户端
	client := NewProxyClient(proxy)

	// 创建请求对象
	req, err := createRequest(method, uri, body, header, cookies)
	// 检查错误
	if err != nil {
		return nil, 0, err
	}

	// 执行请求
	var res *http.Response
	for i := 0; i < RetryTime; i++ {

		PrintRequest(req)
		res, err = client.Do(req)
		if err != nil {
			log.Errorf("[retry %d]http err: %v", i+1, err)
			continue
		}
		PrintRespond(res)

		status = res.StatusCode
		data, err = io.ReadAll(res.Body)
		if err != nil {
			log.Errorf("[retry %d]read body err: %v", i+1, err)
			continue
		}
		// 关闭请求连接
		_ = res.Body.Close()
	}

	return
}

func PrintRequest(r *http.Request) {
	if log.GetLevel() > log.DebugLevel {
		return
	}

	dump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("request: %s", string(dump))
}

func PrintRespond(r *http.Response) {
	if log.GetLevel() > log.DebugLevel {
		return
	}

	dump, err := httputil.DumpResponse(r, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("respond: %s", string(dump))
}

// GetResult 获取远程字节集数据，并返回字节集数据及错误信息
//
// uri 字符串参数，传入请求地址，
// proxy 字符串参数，传入代理地址，
// cookies cookie数组，传入cookie信息。
func GetResult(uri, proxy string, cookies []*http.Cookie) ([]byte, error) {
	// 头部定义
	header := make(map[string]string)
	header["User-Agent"] = UserAgent
	header["referer"] = uri

	// 执行请求
	body, status, err := MakeRequest("GET", uri, proxy, nil, header, cookies)
	// 检查错误
	if err != nil {
		return nil, err
	}

	// 检查状态码
	if http.StatusBadRequest <= status {
		err = fmt.Errorf("%s [Http Status]: %d", uri, status)
	}

	return body, err
}

// GetRoot 获取远程树结构，并返回树结构及错误信息
//
// uri 字符串参数，传入请求地址，
// proxy 字符串参数，传入代理地址，
// cookies cookie数组，传入cookie信息。
func GetRoot(uri, proxy string, cookies []*http.Cookie) (*goquery.Document, error) {
	// 获取远程字节集数据
	data, err := GetResult(uri, proxy, cookies)
	// 检查错误
	if err != nil {
		return nil, err
	}

	// 转换为节点数据
	root, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	// 检查错误
	if err != nil {
		return nil, err
	}

	return root, nil
}

// SavePhoto 下载远程图片到本地，并返回错误信息
//
// uri 字符串参数，远程图片地址，
// savePath 字符串参数，本地保存路径，
// proxy 字符串参数，代理地址，
// needConvert 逻辑参数，是否需要将图片转换为jpg。
func SavePhoto(uri, savePath, proxy string, needConvert bool) error {
	// 创建路径
	err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	// 检查错误
	if err != nil {
		return err
	}

	// 读取远程字节集
	body, err := GetResult(uri, proxy, nil)
	// 检查错误
	if err != nil {
		return err
	}

	// 获取远程图片大小
	length := int64(len(body))
	// 检查大小
	if length == 0 || length < 1024 {
		return fmt.Errorf("远程图片不完整或小于1KB")
	}

	// 保存到本地
	err = saveFile(savePath, body, length)
	// 检查错误
	if err != nil {
		return err
	}

	// 是否需要转换
	if needConvert {
		// 转换为jpg
		err = ConvertJPG(savePath, fmt.Sprintf("%s.jpg", strings.TrimRight(path.Base(savePath), path.Ext(savePath))))
		// 检查
		if err != nil {
			return err
		}

		// 删除源文件
		return os.Remove(savePath)
	}

	return nil
}

// 创建http客户端
func NewProxyClient(proxy string) *http.Client {
	// 初始化
	transport := &http.Transport{
		/* #nosec */
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 如果有代理
	if proxy != "" {
		// 解析代理地址
		proxyURI := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxy)
		}
		// 加入代理
		transport.Proxy = proxyURI
	}

	// 返回客户端
	return &http.Client{
		Transport: transport,
		Timeout:   DefaultTimeOut,
	}
}

// 创建请求对象
func createRequest(method, uri string, body io.Reader, header map[string]string, cookies []*http.Cookie) (*http.Request, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("%s [Request]: %s", uri, err)
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	return req, err
}

// 保存字节集到本地
func saveFile(savePath string, data []byte, length int64) error {
	// 定义错误变量
	var err error

	// 创建路径
	_ = os.MkdirAll(path.Dir(savePath), os.ModePerm)

	// 创建空文件
	f, err := os.Create(savePath)
	// 检查错误
	if err != nil {
		return err
	}

	// 读取数据
	rc := bytes.NewReader(data)
	// 拷贝到指定路径
	_, err = io.Copy(f, rc)
	// 检查错误
	if err != nil {
		return err
	}

	// 关闭连接
	_ = f.Close()

	// 获取文件大小
	local := GetFileSize(savePath)

	// 检查文件一致性
	if length != local {
		// 删除已下载文件
		_ = os.Remove(savePath)
		// 设置错误信息
		err = fmt.Errorf("文件不完成, 下载失败")
	}

	return err
}

type Response struct {
	Respond *http.Response
	Body    []byte
}

func HttpGet(c *http.Client, uri string) (*Response, error) {
	if c == nil {
		c = http.DefaultClient
	}

	rsp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if err := rsp.Body.Close(); err != nil {
		return nil, err
	}

	return &Response{
		Respond: rsp,
		Body:    data,
	}, nil
}
