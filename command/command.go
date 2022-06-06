package command

import (
	"flag"
	"log"
	"s3-gateway/command/vars"
	"s3-gateway/util"
	"strconv"
)

func InitCommand() {
	flag.StringVar(&vars.Port, "Port", "9321", "server port")
	flag.BoolVar(&vars.Debug, "Debug", false, "debug mode")
	flag.BoolVar(&vars.InfoLog, "InfoLog", true, "is print info log")

	var endpointFlags util.ArrayFlags
	var endpointWeight util.ArrayFlags
	flag.Var(&endpointFlags, "Endpoint", "s3 Endpoint")
	flag.Var(&endpointWeight, "EndpointWeight", "s3 Endpoint Weight")
	flag.StringVar(&vars.AccessKey, "AccessKey", "", "s3 AccessKey")
	flag.StringVar(&vars.SecretKey, "SecretKey", "", "s3 SecretKey")

	flag.StringVar(&vars.JWTHeader, "JWTHeader", "Authorization", "JWT Header")
	flag.StringVar(&vars.JWTQuery, "JWTQuery", "jwt", "JWT Query")
	flag.StringVar(&vars.JWTCookie, "JWTCookie", "Authorization", "JWT Cookie")

	var adminList util.ArrayFlags
	flag.Var(&adminList, "AdminList", "Admin List")
	flag.StringVar(&vars.EduBackend, "EduBackend", "", "EduBackend socket")
	flag.StringVar(&vars.Bucket, "Bucket", "", "s3 Bucket")

	flag.Parse()
	if len(endpointWeight) != len(endpointFlags) {
		log.Fatal("The lengths of endpoint and endpoint weight do not match")
	}
	if len(endpointFlags) <= 0 {
		log.Fatal("The lengths of endpoint is zero")
	}
	vars.Endpoint = endpointFlags
	for i, str := range endpointWeight {
		num, err := strconv.Atoi(str)
		if err != nil {
			log.Fatal("The lengths of endpoint weight is not a number")
		}
		for j := 0; j < num; j++ {
			vars.EndpointList = append(vars.EndpointList, vars.Endpoint[i])
		}
	}
	vars.AdminList = adminList
}
