default:
	go mod tidy
	go run main.go
clean:
	rm -rf *.db
