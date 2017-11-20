package config

/**
配置文件格式
*/
type LogConfig struct {
	OutputDir string //输出路径
	ScannerId string //扫描任务唯一ID
	Logfile LogFile
	Task []Task
}

type LogFile struct {
	Dir    string //日志根路径，会递归扫描目录下所有文件
	Order  string //文件排序方式
	Logtype string //日志文件格式，gz
	LogFormat string //日志文件字段格式
	Filter struct {
		Contains    string   //文件包含
		Datepos     []int    //日期起点，日期终点，不配置则不能使用日期过滤
		Datebetween []string //日期过滤
	}
}

type Task struct {
	Taskname   string              //任务名称
	Tasktype   string              //任务类型：count/plain/recall/store
	Url        string              //过滤：包含url
	Groupby    []string            //统计KEY， 按照key计数+1
	Filter     map[string][]string //参数过滤 key:[op,value] op: gt,gte,eq,lte,lt,has,not+op等
	Distanceby []string            //计数排重（参数名称联合排重）任务内计算过的不再计算 - 例如计算UV时按照uid排重
}


/**
任务执行快照，用于断点续传
*/
type Snapshot struct {
	ScannerId    string                    //任务ID
	FinishStatus int                       //任务完成状态
	ScannedCount int                       //已扫描文件数量
	Order        string                    //文件排序方式
	LastLine     uint                      //最终扫描行数
	LastData     map[string]map[string]int //最终任务数据快照
}
