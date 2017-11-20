package logscanner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"scanlog/config"
	"sync"
)

type Context struct {
	Buffer map[string]*DataBuffer // key为任务名称，为每个任务缓存buffer
	EditMutex sync.Mutex      // 同步锁
	Wg sync.WaitGroup
}

type DataBuffer struct {
	UniMapper map[string]bool //去重映射
	CntMapper map[string]int  // 计数映射
	PlainText bytes.Buffer // 文本链接
}

func MakeContext() *Context {
	ctx := new(Context)
	ctx.Buffer = make(map[string]*DataBuffer)
	return ctx
}
func MakeDataBuffer() *DataBuffer {
	buf := new(DataBuffer)
	buf.init()
	return buf
}

func (buf *DataBuffer) init() {
	if buf.UniMapper == nil {
		buf.UniMapper = make(map[string]bool)
	}
	if buf.CntMapper == nil {
		buf.CntMapper = make(map[string]int)
	}
}

/**
保存统计数据到文件
*/
func (ctx *Context) Save2File(logConfig config.LogConfig) {

	base := logConfig.OutputDir

	for task, dataBuf := range ctx.Buffer {

		str := ""
		for key, cnt := range dataBuf.CntMapper {
			str += fmt.Sprintf("%s , %d\n", key, cnt)
		}
		filePath := base + string(os.PathSeparator) + task
		bufStr := bytes.NewBufferString(str)
		ioutil.WriteFile(filePath, bufStr.Bytes(), os.ModeAppend)

		fmt.Println("已保存文件至:", filePath)
	}

}

func (ctx *Context) Save2FileTemp(logConfig config.LogConfig) {

	base := logConfig.OutputDir

	fmt.Println(ctx.Buffer)
	for task, dataBuf := range ctx.Buffer {

		str := ""
		for key, cnt := range dataBuf.CntMapper {
			str += fmt.Sprintf("%s , %d\n", key, cnt)
		}
		filePath := base + string(os.PathSeparator) + task+"-temp"
		bufStr := bytes.NewBufferString(str)
		ioutil.WriteFile(filePath, bufStr.Bytes(), os.ModeAppend)

		fmt.Println("已保存文件至:", filePath)
	}

}
func createIfNotExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		os.Create(path)
		return false, nil
	}
	return false, err
}



/*
	数据计数
*/
func (ctx *Context) CntData(task, key, unikey string) {
	var dataBuf *DataBuffer
	ctx.EditMutex.Lock()

	//初始化
	if ctx.Buffer[task] != nil {
		dataBuf = ctx.Buffer[task]
	} else {
		dataBuf = MakeDataBuffer()
	}

	if !dataBuf.UniMapper[unikey] || len(unikey)==0 {
		dataBuf.UniMapper[unikey] = len(unikey)!=0
		dataBuf.CntMapper[key] = dataBuf.CntMapper[key] + 1
		ctx.Buffer[task] = dataBuf
	}
	ctx.EditMutex.Unlock()
}
