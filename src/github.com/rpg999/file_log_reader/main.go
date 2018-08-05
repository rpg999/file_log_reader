package build

import (
	"flag"
	"log"
	"strings"
	"os"
	"github.com/rpg999/file_log_reader"
)

var fileList = flag.String("file_list", "/tmp/test_files/first_format.log|/tmp/test_files/second_format.log", "List of files to track, separated by `|`")
var logFormat = flag.String("log_format", "first_format", "Type of logs in the file")

func main()  {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(*fileList) == 0 {
		log.Fatal("no file name is specified")
	}

	var paths []string
	for _, path := range strings.Split(*fileList, "|") {
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		}
	}

	if len(paths) == 0 {
		log.Fatal("no valid files for tracking")
	}

	if logStream, err := file_log_reader.TrackFiles(paths); err != nil {
		for l := range logStream {
			log.Println(l)
		}
	}

	log.Println("End of program")
}