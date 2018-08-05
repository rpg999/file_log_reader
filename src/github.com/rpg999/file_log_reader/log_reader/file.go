package log_reader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type LogFile struct {
	FullPath string `json:"full_path"`
	Format   string `json:"format"`
}

func (l *LogFile) UnmarshalJSON(b []byte) error {
	type plain LogFile
	if err := json.Unmarshal(b, (*plain)(l)); err != nil {
		return err
	}

	if len(l.FullPath) == 0{
		return fmt.Errorf("path to file is empty")
	}

	if l.Format != FirstFormat && l.Format != SecondFormat {
		return fmt.Errorf("invalid log format: %s", l.Format)
	}

	return nil
}

func ParseFileList(path string) ([]LogFile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading file: %s", err)
	}

	var logFiles = &struct {
		Files []LogFile `json:"files"`
	}{}

	err = json.Unmarshal(b, logFiles)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling file list: %s", err)
	}

	return logFiles.Files, nil
}
