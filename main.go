package main

import router "github.com/bncunha/erp-api/src/api"

func main() {
	r := router.NewRouter()
	r.SetupRoutes()
	r.Start()
}