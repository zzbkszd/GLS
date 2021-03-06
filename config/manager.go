package config

import (
	"io/ioutil"
	"encoding/json"
	"strings"
	"fmt"
)
/**
加载配置文件
 */
func (config *LogConfig) Load(f string){

	str,readerr:= ioutil.ReadFile(f)
	if readerr != nil {
		panic(readerr)
	}
	jsonerr := json.Unmarshal(str,config)
	if jsonerr!=nil {
		panic(jsonerr)
	}

	fmt.Println("Load config file from ",f," success!")

}

/**
过滤文件名称
 */
func (config *LogConfig) FilterFile (name string) bool{

	if ! (strings.Contains(name,config.Logfile.Filter.Contains) && strings.HasSuffix(name,config.Logfile.Logtype)){
		fmt.Println("filt file fail: ",name," by contains")
		return false
	}

	if len(config.Logfile.Filter.Datepos)>=2 {
		start := config.Logfile.Filter.Datepos[0]
		end := config.Logfile.Filter.Datepos[1]
		date := string([]rune(name)[start+1:end])
		if len(config.Logfile.Filter.Datebetween) >=2 {
			if date < config.Logfile.Filter.Datebetween[0] || date > config.Logfile.Filter.Datebetween[1]{
				fmt.Println("filt file fail: ",name," by date:",date)
				return false
			}
		}
	}
	return true
}

/**
输出配置json
 */
func (config LogConfig) String() string {
	res := ""
	json,err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	res += string(json)
	return res
}