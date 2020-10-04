default:
	go mod tidy
	go build
	./photo-gallery
clean:
	rm -rf *.db
