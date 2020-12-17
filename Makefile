build: 
	go build -o cmd/aphrodite-web.runtime -ldflags "-s -w" startup.go