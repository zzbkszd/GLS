package logscanner

import (
	"fmt"
	"io/ioutil"
	"os"
	"scanlog/config"
)

func Startup(config config.LogConfig) {
	//TODO 通过对文件做过滤和稳定排序以支持断点续传
	files := ListDir(config.Logfile.Dir)
	ctx := MakeContext()
	skip := 0 //跳过指定个数的文件
	for _,f:=range files {
		fmt.Println(f)
		if config.FilterFile(f){
			if(skip>0){
				skip--
				fmt.Println("skip file "+f)
				continue
			}
			ls := LogScanner{path:f,config:config,ctx:ctx}
			ls.Scan()
		}
	}
	fmt.Println("read log finish!")
	ctx.Save2File(config)
}


//列出目录下的所有文件
func ListDir (base string) []string {
	fileList := make ([]string,0)
	dir,err := ioutil.ReadDir(base)
	if err != nil {
		panic(err)
		return fileList
	}
	PthSep := string(os.PathSeparator)
	for _,fi := range dir {
		if fi.IsDir(){
			subs := ListDir(base+PthSep+fi.Name())
			for _,s := range subs {
				fileList = append(fileList,s)
			}
		}else {
			fileList = append(fileList, base+PthSep+fi.Name())
		}
	}
	return fileList

}