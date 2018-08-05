package log_reader

import (
	"errors"
	"fmt"
	"github.com/hpcloud/tail"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	FirstFormat  = "first_format"
	SecondFormat = "second_format"
)

type Log struct {
	// Time of recording from the log file
	LogTime time.Time `bson:"log_time"`

	// Log message from the log file
	LogMsg string `bson:"log_msg"`

	// The path to the file from which the message was received
	FileName string `bson:"file_name"`

	// • log_format (String) - log format (first_format | second_format)
	LogFormat string `bson:"log_format"`
}

type Result struct {
	Log *Log
	Err error
}

// TrackFiles tracks files and convert the result to the log object and push it to log stream
func TrackFiles(done chan struct{}, files []LogFile) (<-chan Result, error) {
	logStream, err := trackFiles(done, files)
	if err != nil {
		return nil, err
	}

	return logStream, nil
}

var ErrNoLogStreams = errors.New("no log streams are available")

// trackFiles range over files and start tracking file
func trackFiles(done chan struct{}, files []LogFile) (<-chan Result, error) {
	var logStreams []<-chan Result
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
func trackFile(done chan struct{}, file LogFile) (<-chan Result, error) {
	t, err := tail.TailFile(file.FullPath, tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}

	var parseLine func(string, string) Result
	if file.Format == FirstFormat {
		parseLine = parseLineFirstFormat
	} else {
		parseLine = parseLineSecondFormat
	}

	logStream := make(chan Result)
	go func() {
		defer t.Stop()
		defer close(logStream)

		for {
			select {
			case <-done:
				return
			case line := <-t.Lines:
				logStream <- parseLine(line.Text, file.FullPath)
			}
		}
	}()

	return logStream, nil
}

// fanIn combines data from multiple log channels to one
func fanIn(done chan struct{}, channels ...<-chan Result) <-chan Result {
	var wg sync.WaitGroup
	multiplexedStream := make(chan Result)

	multiplex := func(c <-chan Result) {
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

// parseLineFirstFormat parse line of the first format and return the result
func parseLineFirstFormat(line, path string) Result {
	r := Result{}
	parts := strings.Split(line, " | ")
	if len(parts) != 2 {
		r.Err = fmt.Errorf("%s: unable to parse text string: %s", FirstFormat, line)
		return r
	}

	// Feb 1, 2018 at 3:04:05pm (UTC)
	t, err := time.Parse("Jan 2, 2006 at 3:04:05pm (MST)", parts[0])
	if err != nil {
		r.Err = fmt.Errorf("%s: unable to parse string to time: %s", FirstFormat, line)
		return r
	}

	l := &Log{
		LogMsg:    parts[1],
		FileName:  path,
		LogFormat: FirstFormat,
		LogTime:   t,
	}
	r.Log = l

	return r
}

// parseLineSecondFormat parse line of the first format and return the result
func parseLineSecondFormat(line, path string) Result {
	r := Result{}
	parts := strings.Split(line, " | ")
	if len(parts) != 2 {
		r.Err = fmt.Errorf("%s: unable to parse text string: %s", SecondFormat, line)
		return r
	}

	// 2018-02-01T15:04:05Z
	t, err := time.Parse("2006-02-01T15:04:05Z", parts[0])
	if err != nil {
		r.Err = fmt.Errorf("%s: unable to parse string to time: %s", FirstFormat, line)
		return r
	}

	l := &Log{
		LogMsg:    parts[1],
		FileName:  path,
		LogFormat: SecondFormat,
		LogTime:   t,
	}
	r.Log = l

	return r
}
