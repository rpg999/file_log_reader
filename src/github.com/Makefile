build:
	gb build
run: build
	bin/file_log_reader
init:
	go get -u github.com/constabulary/gb/...
	gb vendor fetch github.com/mongodb/mongo-go-driver/mongo
	gb vendor fetch github.com/hpcloud/tail