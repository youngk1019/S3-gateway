package main

import (
	"s3-gateway/command"
	"s3-gateway/log"
	"s3-gateway/routers"
)

func main() {
	command.InitCommand()
	log.InitLogger()
	routers.InitRouter()
	routers.Run()
}
