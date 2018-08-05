package reader

import (
	"errors"
	"github.com/hpcloud/tail"
	"log"
	"sync"
	"time"
	"strings"
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

type Result struct {
	Log *Log
	Err error
}

// TrackFiles tracks files and convert the result to the log object and push it to log stream
func TrackFiles(files []LogFile) (<-chan *Log, error) {
	// Done channel can make all goroutines exist safely
	done := make(chan struct{})
	logStream, err := trackFiles(done, files)
	if err != nil {
		close(done)
		return nil, err
	}

	return logStream, nil
}

var ErrNoLogStreams = errors.New("no log streams are available")

// trackFiles range over files and start tracking file
func trackFiles(done chan struct{}, files []LogFile) (<-chan *Log, error) {
	var logStreams []<-chan *Log
	for _, file := range files {
		if logStream, err := trackFile(done, file); err != nil {
			log.Printf("error on tracking file %s: %s", file, err)
		} else {
			logStreams = append(logStreams, logStream)
		}
	}

	if len(logStreams) == 0 {
		return nil, ErrNoLogStreams
	}

	return fanIn(done, logStreams...), nil
}

// Starts tracking the file and sends back the channel where the log will be sent
func trackFile(done chan struct{}, file LogFile) (<-chan *Log, error) {
	t, err := tail.TailFile(file.FullPath, tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}

	logStream := make(chan *Log)
	go func() {
		defer t.Stop()
		defer close(logStream)

		for {
			select {
			case <-done:
				return
			case line := <-t.Lines:
				parts := strings.Split(line.Text, " | ")
				if len(parts) != 2 {
					log.Printf("unable to parse text string: %s", line.Text)
					continue
				}

				time.Parse("", parts[0])

				l := &Log{

					LogMsg: parts[1],
					FileName: file.FullPath,
					LogFormat: file.Format,
				}
				logStream <- l
			}
		}
	}()

	return logStream, nil
}

// fanIn combines data from multiple log channels to one
func fanIn(done chan struct{}, channels ...<-chan *Log) <-chan *Log {
	var wg sync.WaitGroup
	multiplexedStream := make(chan *Log)

	multiplex := func(c <-chan *Log) {
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