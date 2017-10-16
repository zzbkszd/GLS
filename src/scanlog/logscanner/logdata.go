package logscanner

import (
	"strconv"
	"strings"
	"fmt"
	"time"
)

type LogData struct {
	ClientIp     string            //客户端ip
	Date         string            //请求日期-YMD
	Timestamp    int64             //请求时间戳
	Method       string            //请求方式
	RequesStr    string            //请求参数串
	RequestUri   string            //请求接口uri
	ResponseCode int               //返回码
	Host         string            //服务
	BodyLength   int               //返回数据长度
	CostTime     float64           //处理请求耗时
	Params       map[string]string //参数列表
	UserAgent    string            //用户UA
}

/**
将日志按字段划分
TODO 应该采取正则提取的形式，利于支持多种样式的日志
*/
func parselog(line string) []string {
	fields := strings.Fields(line)
	result := make([]string, 0, 12)

	mark := false
	var cache = string("")
	//合并用引号括起来的包含空格的字段
	for _, f := range fields {
		if strings.HasPrefix(f, "\"") {
			mark = true
		}
		if strings.HasSuffix(f, "\"") {
			mark = false
		}
		cache += " "
		cache += f
		if !mark {
			result = append(result, cache)
			cache = string("")
		}
	}
	return result

}

/**
解析日志，构造结构体
TODO 需要支持按照输入格式来动态分析字段
TODO 支持调用插件以支持更多样式的日志分析
*/
func makeData(logstr string) LogData {
	log := parselog(logstr)
	//fmt.Println(log)
	if len(log) != 13 {
		return LogData{}
	}
	//fmt.Println(log)
	date,timestamp := parseDate(string([]rune(log[3])[2:len(log[3])]),"nginx_ori")
	requestStr := log[5]
	requestFields := strings.Fields(requestStr)
	if len(requestFields) < 2 {
		return LogData{}
	}
	method := string([]rune(requestFields[0])[2:len(requestFields[0])])
	splitUrl := strings.Split(requestFields[1], "?")
	requrl := ""
	if len(splitUrl) == 2 {
		requrl = splitUrl[1]
	}
	responseCode, _ := strconv.Atoi(log[6])
	bodyLength, _ := strconv.Atoi(log[7])
	costtime, _ := strconv.ParseFloat(log[11], 32)
	//post上传的参数
	postData := log[12]
	if len(postData)!=0 && postData!="-" {
		requrl+="&"+postData
	}
	params := parseParam(requrl)

	data := LogData{ClientIp: log[0], Date: date, Method: method, RequesStr: requrl, ResponseCode: responseCode,
		Host: log[10], UserAgent: log[9], CostTime: costtime, BodyLength: bodyLength, Params: params , RequestUri: splitUrl[0],Timestamp:timestamp}
	return data
}

func parseDate(date,dateSrc string ) (string,int64) {
	var day,month,year,hour,min,sec string
	var strDate, strTime string
	var timestamp int64
	if dateSrc == "nginx_ori" {
		fmt.Sscanf(date,"%2s/%3s/%4s:%2s:%2s:%2s",&day,&month,&year,&hour,&min,&sec)
		switch month {
			case "Jan": month="01"
			case "Feb": month="02"
			case "Mar": month="03"
			case "Apr": month="04"
			case "May": month="05"
			case "Jun": month="06"
			case "Jul": month="07"
			case "Aug": month="08"
			case "Sep": month="09"
			case "Oct": month="10"
			case "Nov": month="11"
			case "Dec": month="12"
		}
		strDate = fmt.Sprintf("%s-%s-%s",year,month,day)
		strTime = fmt.Sprintf("%s:%s:%s",hour,min,sec)
		unixTime,_ := time.Parse("2006-01-02 15:04:05", strDate+" "+strTime)
		//fmt.Println(unixTime.Unix())
		timestamp = unixTime.Unix()
		//if err != nil {
		//	fmt.Println(timestamp)
		//}

	}
	return strDate,timestamp
}

/**
解析get请求串中的参数
*/
func parseParam(param string) map[string]string {
	mapper := make(map[string]string)
	pairs := strings.Split(param, "&")

	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		v := ""
		if len(kv) == 2 {
			v = kv[1]
		}
		mapper[kv[0]] = v
	}
	return mapper

}
