package logscanner

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strings"
	"scanlog/config"
	"strconv"
	"regexp"
	"time"
)

type LogScanner struct {
	path string // 文件路径
	name string // 协程名称
	config config.LogConfig //配置信息
	ctx *Context //上下文信息
	endchan chan int
}

/**
日志处理工具
*/
func (ls *LogScanner) Scan() {
	fmt.Println("start read log file ", ls.path)
	reg := regexp.MustCompile("^(?P<ip>[0-9\\.]{7,15}) - - \\[(?P<date>[^\\s]+) \\+[0-9]{4}\\] \"(?P<method>[A-Z]{0,4}) (?P<uri>[^\\s]+) [^\"]+\" (?P<code>[0-9]{3}) (?P<bodylength>[0-9]*) \".*?\" \"(?P<useragent>[^\"]*)\" \"(?P<host>[^\"]*?)\" (?P<costtime>[0-9.]*) \"(?P<postdata>[^\"]*)\"$")
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
	start := time.Now().UnixNano()
	for {
		//读取数据
		line, prefix, _ := bufreader.ReadLine()
		logStr := string(line)
		for prefix {
			l,pf,_ := bufreader.ReadLine()
			logStr = logStr+string(l)
			prefix = pf
		}
		if line == nil {
			break
		}

		//计数
		readLen+=(int64)(len(logStr))
		c++

		log := makeData(logStr,reg)

		//遍历所有task
		for i:=0;i<len(ls.config.Task);i++ {
			task := ls.config.Task[i]
			//TODO 条件过滤
			if filterLog(log,task) {
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
			fmt.Println("make cost:",time.Now().UnixNano()-start)
			start = time.Now().UnixNano()
			fmt.Printf("Task %s process : %.2f%%\n",ls.name,float64(float64(readLen)/float64(totalLen*5))*100)
		}
	}
	ls.ctx.Wg.Done()
	ls.endchan <- 1
}

func filterLog( log *LogData, task config.Task) bool{
	uri := task.Url
	res := true
	res = res && strings.Contains(log.RequestUri, uri)
	if !res {
		return res
	}
	filter := task.Filter
	for k,v := range filter {
		lk := log.Params[k]
		op := v[0]
		switch op {
		case "eq":
			res = res && (lk == v[1])
		case "neq":
			res = res && (lk != v[1])
		}
	}
	return res

}


func genKey(log *LogData,fields []string) string{
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


