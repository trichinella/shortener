package main

func main() {
	server := CreateServer(NewConfig())
	server.Start()
}
