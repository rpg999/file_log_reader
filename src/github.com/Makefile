build:
	gb build
test-run: build
	bin/file_log_reader -file_list=src/github.com/rpg999/file_log_reader/test_files/test_file_list.json
init:
	go get -u github.com/constabulary/gb/...
	gb vendor fetch github.com/mongodb/mongo-go-driver/mongo
	gb vendor fetch github.com/hpcloud/tail