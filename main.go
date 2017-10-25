package main

import (
	"scanlog/logscanner"
	"scanlog/config"
	"flag"
)

func main(){
	//configPath := flag.String("config","logconfig.json","Config file path.")
	configPath := flag.String("config","src\\scanlog\\logconfig.json","Config file path.")
	//autosave := flag.Uint("autosave",100000,"N uint = create snapshot every N lines.")
	config := config.LogConfig{}
	config.Load(*configPath)
	logscanner.Startup(config)
}


