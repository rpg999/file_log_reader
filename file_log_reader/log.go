package file_log_reader

import (
	"time"
	"sync"
)

type Log struct {
	// Time of recording from the log file
	LogTime time.Time `json:"log_time"`

	// Log message from the log file
	LogMsg string `json:"log_msg"`

	// The path to the file from which the message was received
	FileName string `json:"file_name"`

	// • log_format (String) - log format (first_format | second_format)
	LogFormat string `json:"log_format"`
}

// trackFiles tracks files of paths and convert the result to the log object and push it to log stream
func TrackFiles(paths []string) <-chan Log {
	done := make(chan struct {})

	pathRoutines := make([]<-chan Log, len(paths))
	trackFiles(done, paths)

	return fanIn(done, pathRoutines...)
}

func trackFiles(done chan struct{}, paths []string) {

}

// fanIn combines data from multiple channels
func fanIn(done chan struct{}, channels ...<-chan Log) <-chan Log {
	var wg sync.WaitGroup
	multiplexedStream := make(chan Log)

	multiplex := func(c <- chan Log) {
		defer wg.Done()
		for i := range c {
			select {
				case <-done:
				return

				case multiplexedStream <- i:
			}
		}
	}

	// Select from all channels
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}