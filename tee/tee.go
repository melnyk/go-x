package tee

import (
	"sync"
	"sync/atomic"

	"go.melnyk.org/mlog"
)

type teelogger struct {
	loggers     []mlog.Logger
	name        string
	customlevel bool
	level       mlog.Level
}

// Interface implementation check
var (
	_ mlog.Logger = &teelogger{}
)

func (l *teelogger) Verbose(msg string) {
	if l.level < mlog.Verbose {
		return
	}

	for _, logger := range l.loggers {
		logger.Verbose(msg)
	}
}

func (l *teelogger) Info(msg string) {
	if l.level < mlog.Info {
		return
	}

	for _, logger := range l.loggers {
		logger.Info(msg)
	}
}

func (l *teelogger) Warning(msg string) {
	if l.level < mlog.Warning {
		return
	}

	for _, logger := range l.loggers {
		logger.Warning(msg)
	}
}

func (l *teelogger) Error(msg string) {
	if l.level < mlog.Error {
		return
	}

	for _, logger := range l.loggers {
		logger.Error(msg)
	}
}

func (l *teelogger) Panic(msg string) {
	if l.level < mlog.Panic {
		return
	}

	for _, logger := range l.loggers {
		logger.Panic(msg)
	}
}

func (l *teelogger) Fatal(msg string) {
	for _, logger := range l.loggers {
		logger.Fatal(msg)
	}
}

func (l *teelogger) Event(level mlog.Level, cb func(evt mlog.Event)) {
	if l.level < level {
		return
	}

	for _, logger := range l.loggers {
		logger.Event(level, cb)
	}
}

type teelogbook struct {
	logbooks     []mlog.Logbook
	defaultlevel mlog.Level
	mu           sync.Mutex
	loggers      map[string]*teelogger
}

// Interface implementation check
var (
	_ mlog.Logbook = &teelogbook{}
)

func (lb *teelogbook) SetLevel(name string, level mlog.Level) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Set default level...
	if name == mlog.Default {
		lb.defaultlevel = level
		// so update all non-custom loggers too
		for _, v := range lb.loggers {
			if !v.customlevel {
				atomic.StoreUint32((*uint32)(&v.level), uint32(level))
			}
		}
	}

	// Set level for dedicated logger makes it custom
	if l, ok := lb.loggers[name]; ok {
		atomic.StoreUint32((*uint32)(&l.level), uint32(level))
		l.customlevel = true
		return nil
	}

	// Well...logger name is unknown, so add new one to logbook
	loggers := make([]mlog.Logger, 0)
	for _, v := range lb.logbooks {
		loggers = append(loggers, v.Joiner().Join(name))
	}

	l := &teelogger{
		name:        name,
		level:       level,
		customlevel: true,
		loggers:     loggers,
	}

	lb.loggers[name] = l

	return nil
}

func (lb *teelogbook) Levels() mlog.Levels {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lvs := make(mlog.Levels)

	lvs[mlog.Default] = lb.defaultlevel
	for k, v := range lb.loggers {
		lvs[k] = mlog.Level(atomic.LoadUint32((*uint32)(&v.level)))
	}

	return lvs
}

func (lb *teelogbook) Joiner() mlog.Joiner {
	// logbook provides this interface in this implemenation
	return lb
}

func (lb *teelogbook) Join(name string) mlog.Logger {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Return exist logger
	if l, ok := lb.loggers[name]; ok {
		return l
	}

	// ... or create new one
	loggers := make([]mlog.Logger, 0)
	for _, v := range lb.logbooks {
		loggers = append(loggers, v.Joiner().Join(name))
	}

	l := &teelogger{
		name:        name,
		level:       lb.defaultlevel,
		customlevel: false,
		loggers:     loggers,
	}

	lb.loggers[name] = l

	return l
}

// NewLogbook returns interface to logbook implementation
func NewLogbook(logbooks ...mlog.Logbook) mlog.Logbook {
	lbs := make([]mlog.Logbook, 0)
	for _, v := range logbooks {
		// Control is moved to tee logbook. It can be configured after this
		v.SetLevel(mlog.Default, mlog.Verbose)
		lbs = append(lbs, v)
	}

	lb := &teelogbook{
		logbooks: lbs,
		loggers:  make(map[string]*teelogger),
	}

	return lb
}
