build:
	go mod tidy
	go build
	./photo-gallery
run:
	go mod tidy
	go run *.go
clean:
	rm -rf *.db
default:
	run
