package main

import "api/routers"

func main() {
	r := routers.Routes()
	r.Run()
}