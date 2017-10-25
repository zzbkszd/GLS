package logscanner

import (
	"fmt"
	"io/ioutil"
	"os"
	"scanlog/config"
	"runtime"
	"strconv"
)

func Startup(config config.LogConfig) {
	maxProcs := runtime.NumCPU()   //获取cpu个数
	runtime.GOMAXPROCS(maxProcs)  //限制同时运行的goroutines数量
	//TODO 通过对文件做过滤和稳定排序以支持断点续传

	files := ListDir(config.Logfile.Dir,config)
	ctx := MakeContext()
	i := 0 //跳过指定个数的文件
	gocnt := 0
	end := make(chan int )

	for ;i<len(files); {
		if gocnt < maxProcs {
			ls := LogScanner{path:files[i],config:config,ctx:ctx,name:files[i]+"-go-"+strconv.Itoa(i),endchan:end}
			ctx.Wg.Add(1)
			go ls.Scan()
			gocnt++
			i++
		}else {
			select {
				case <-end :
					gocnt--
			}
		}
	}
	//for _,f:=range files {
	//	if(skip>0){
	//		skip--
	//		fmt.Println("skip file "+f)
	//		continue
	//	}
	//	ls := LogScanner{path:f,config:config,ctx:ctx,name:f+"-go-"+strconv.Itoa(gocnt)}
	//	gocnt++
	//	ctx.Wg.Add(1)
	//	//TODO - 用线程池来处理
	//	go ls.Scan()
	//}
	ctx.Wg.Wait()
	fmt.Println("read log finish!")
	ctx.Save2File(config)
}


//列出目录下的所有文件
func ListDir (base string, config config.LogConfig) []string {
	fileList := make ([]string,0)
	dir,err := ioutil.ReadDir(base)
	if err != nil {
		panic(err)
		return fileList
	}
	PthSep := string(os.PathSeparator)
	for _,fi := range dir {
		fullName := base+PthSep+fi.Name()
		if fi.IsDir(){
			subs := ListDir(fullName,config)
			for _,s := range subs {
				fileList = append(fileList,s)
			}
		}else {
			if config.FilterFile(fi.Name()) {
				fileList = append(fileList, fullName)
			}
		}
	}
	return fileList

}