package main

func main() {
	server := NewAPIServer(":8443", "./cert/server.crt", "./cert/server.key")
	server.Run()
}
