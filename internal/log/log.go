package log

import (
	"sync"
)

type Log struct {
	mu sync.Mutex

	Dir    string
	Config Config

	activeSegment *segment
	segments      []*segment
}

func NewLog(dir string, c Config) (*Log, error) {

}
