
build: clean
	mkdir build
	go build -o build/s3ftp

clean:
	rm -rf build/
