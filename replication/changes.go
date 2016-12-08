package replication

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

func NewChangesetDownloader(dest, url string, seq int, interval time.Duration) *downloader {
	dl := newDownloader(dest, url, seq, interval)
	dl.fileExt = ".osm.gz"
	dl.stateExt = ".state.txt"
	dl.stateTime = parseYamlTime
	go dl.fetchNextLoop()
	return dl
}

type changesetState struct {
	Time     yamlStateTime `yaml:"last_run"`
	Sequence int           `yaml:"sequence"`
}

type yamlStateTime struct {
	time.Time
}

func (y *yamlStateTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ts string
	if err := unmarshal(&ts); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -07:00", ts)
	y.Time = t
	return err
}

func parseYamlStateFile(filename string) (changesetState, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return changesetState{}, err
	}
	return parseYamlState(b)
}

func parseYamlState(b []byte) (changesetState, error) {
	state := changesetState{}
	if err := yaml.Unmarshal(b, &state); err != nil {
		return changesetState{}, err
	}
	return state, nil
}

func parseYamlTime(filename string) (time.Time, error) {
	state, err := parseYamlStateFile(filename)
	if err != nil {
		return time.Time{}, err
	}
	return state.Time.Time, nil
}
