package logscanner

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strings"
	"scanlog/config"
	"strconv"
	"bytes"
)

type LogScanner struct {
	path string // 文件路径
	config config.LogConfig //配置信息
	ctx *Context //上下文信息
}


/**
日志处理工具
*/
func (ls *LogScanner) Scan() {

	fmt.Println("start read log file ", ls.path)
	file, e := os.Open(ls.path)
	if e != nil {
		panic(e)
		return
	}
	defer file.Close()

	reader, re := gzip.NewReader(file)
	if re != nil {
		panic(re)
		return
	}

	//文件输入
	bufreader := bufio.NewReader(reader)

	c := 0
	stat,_ := file.Stat()
	var totalLen = stat.Size()
	var readLen int64
	readLen=0
	for {
		//读取数据
		line, prefix, _ := bufreader.ReadLine()
		byteBuf := bytes.NewBuffer(line)
		for prefix {
			l,pf,_ := bufreader.ReadLine()
			byteBuf.Write(l)
			prefix = pf
		}
		line = byteBuf.Bytes()
		if line == nil {
			break
		}

		//计数
		readLen+=(int64)(len(string(line)))
		c++

		log := makeData(string(line))

		//遍历所有task
		for i:=0;i<len(ls.config.Task);i++ {
			task := ls.config.Task[i]
			uri := task.Url
			//TODO 条件过滤
			if strings.Contains(log.RequestUri, uri) {
				//计算key与unikey
				key := genKey(log,task.Groupby)
				var unikey string
				if len(task.Distanceby)>0 {
					unikey = key + genKey(log,task.Distanceby)
				}
				ls.ctx.CntData(task.Taskname,key,unikey)
				//TODO 支持调用插件以扩展功能
			}
		}

		//TODO 输出进度
		if c%100000 == 0 && c != 0 {
			fmt.Printf("%.2f%%\n",float64(float64(readLen)/float64(totalLen*5))*100)
		}
	}
	//TODO 保存到文件
	//UNDO 留在context统一处理

}

func genKey(log LogData,fields []string) string{
	key := ""
	for _,v := range fields {
		switch v {
		case "&clientIp":
			key+=log.ClientIp
		case "&date":
			key+=log.Date
		case "&method":
			key+=log.Method
		case "&responseCode":
			key+=strconv.Itoa(log.ResponseCode)
		case "&host":
			key+=log.Host
		case "&userAgent":
			key+=log.UserAgent
		default:
			key+=log.Params[v]
		}
		key+="_"
	}
	return key
}


