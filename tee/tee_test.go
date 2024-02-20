package tee

import (
	"bytes"
	"strings"
	"testing"

	"go.melnyk.org/mlog"
	"go.melnyk.org/mlog/console"
	"go.melnyk.org/mlog/nolog"
)

func utilFillLog(log mlog.Logger) {
	log.Verbose("verbose msg")
	log.Info("info msg")
	log.Warning("warning msg")
	log.Error("error msg")
	log.Panic("panic msg")
	log.Fatal("fatal msg")
}

func utilFillLogCB(log mlog.Logger) {
	log.Event(mlog.Verbose, func(e mlog.Event) {
		e.String("msg", "verbose msg")
	})
	log.Event(mlog.Info, func(e mlog.Event) {
		e.String("msg", "info msg")
	})
	log.Event(mlog.Warning, func(e mlog.Event) {
		e.String("msg", "warning msg")
	})
	log.Event(mlog.Error, func(e mlog.Event) {
		e.String("msg", "error msg")
	})
	log.Event(mlog.Panic, func(e mlog.Event) {
		e.String("msg", "panic msg")
	})
	log.Event(mlog.Fatal, func(e mlog.Event) {
		e.String("msg", "fatal msg")
	})
}

func utilTestLevel(t *testing.T, out string, level mlog.Level) {
	if inc := strings.Contains(out, "FATAL"); inc != (level >= mlog.Fatal) {
		t.Fatal("Failed FATAL level check", out)
	}

	if inc := strings.Contains(out, "fatal msg"); inc != (level >= mlog.Fatal) {
		t.Fatal("Failed FATAL level check (msg output)")
	}

	if inc := strings.Contains(out, "PANIC"); inc != (level >= mlog.Panic) {
		t.Fatal("Failed PANIC level check", out)
	}

	if inc := strings.Contains(out, "panic msg"); inc != (level >= mlog.Panic) {
		t.Fatal("Failed PANIC level check (msg output)")
	}

	if inc := strings.Contains(out, "ERROR"); inc != (level >= mlog.Error) {
		t.Fatal("Failed ERROR level check", out)
	}

	if inc := strings.Contains(out, "error msg"); inc != (level >= mlog.Error) {
		t.Fatal("Failed ERROR level check (msg output)")
	}

	if inc := strings.Contains(out, "WARNING"); inc != (level >= mlog.Warning) {
		t.Fatal("Failed WARNING level check")
	}

	if inc := strings.Contains(out, "warning msg"); inc != (level >= mlog.Warning) {
		t.Fatal("Failed WARNING level check (msg output)")
	}

	if inc := strings.Contains(out, "INFO"); inc != (level >= mlog.Info) {
		t.Fatal("Failed INFO level check")
	}

	if inc := strings.Contains(out, "info msg"); inc != (level >= mlog.Info) {
		t.Fatal("Failed INFO level check (msg output)")
	}

	if inc := strings.Contains(out, "VERBOSE"); inc != (level >= mlog.Verbose) {
		t.Fatal("Failed VERBOSE level check")
	}

	if inc := strings.Contains(out, "verbose msg"); inc != (level >= mlog.Verbose) {
		t.Fatal("Failed VERBOSE level check (msg output)")
	}
}

func TestLoggerFatal(t *testing.T) {
	logbookA := nolog.NewLogbook()
	buf := &bytes.Buffer{}
	logbookB := console.NewLogbook(buf)
	logbook := NewLogbook(logbookA, logbookB)

	logger := logbook.Joiner().Join("test")

	logbook.SetLevel("test", mlog.Fatal)
	utilFillLog(logger)
	utilTestLevel(t, buf.String(), mlog.Fatal)

	buf.Reset()
	utilFillLogCB(logger)
	utilTestLevel(t, buf.String(), mlog.Fatal)
}

func TestLoggerPanic(t *testing.T) {
	logbookA := nolog.NewLogbook()
	buf := &bytes.Buffer{}
	logbookB := console.NewLogbook(buf)
	logbook := NewLogbook(logbookA, logbookB)

	logger := logbook.Joiner().Join("test")

	logbook.SetLevel("test", mlog.Panic)
	utilFillLog(logger)
	utilTestLevel(t, buf.String(), mlog.Panic)

	buf.Reset()
	utilFillLogCB(logger)
	utilTestLevel(t, buf.String(), mlog.Panic)
}

func TestLoggerError(t *testing.T) {
	logbookA := nolog.NewLogbook()
	buf := &bytes.Buffer{}
	logbookB := console.NewLogbook(buf)
	logbook := NewLogbook(logbookA, logbookB)

	logger := logbook.Joiner().Join("test")

	logbook.SetLevel("test", mlog.Error)
	utilFillLog(logger)
	utilTestLevel(t, buf.String(), mlog.Error)

	buf.Reset()
	utilFillLogCB(logger)
	utilTestLevel(t, buf.String(), mlog.Error)
}

func TestLoggerWarning(t *testing.T) {
	logbookA := nolog.NewLogbook()
	buf := &bytes.Buffer{}
	logbookB := console.NewLogbook(buf)
	logbook := NewLogbook(logbookA, logbookB)

	logger := logbook.Joiner().Join("test")

	logbook.SetLevel("test", mlog.Warning)
	utilFillLog(logger)
	utilTestLevel(t, buf.String(), mlog.Warning)

	buf.Reset()
	utilFillLogCB(logger)
	utilTestLevel(t, buf.String(), mlog.Warning)
}

func TestLoggerInfo(t *testing.T) {
	logbookA := nolog.NewLogbook()
	buf := &bytes.Buffer{}
	logbookB := console.NewLogbook(buf)
	logbook := NewLogbook(logbookA, logbookB)

	logger := logbook.Joiner().Join("test")

	logbook.SetLevel("test", mlog.Info)
	utilFillLog(logger)
	utilTestLevel(t, buf.String(), mlog.Info)

	buf.Reset()
	utilFillLogCB(logger)
	utilTestLevel(t, buf.String(), mlog.Info)
}

func TestLoggerVerbose(t *testing.T) {
	logbookA := nolog.NewLogbook()
	buf := &bytes.Buffer{}
	logbookB := console.NewLogbook(buf)
	logbook := NewLogbook(logbookA, logbookB)

	logger := logbook.Joiner().Join("test")

	logbook.SetLevel("test", mlog.Verbose)
	utilFillLog(logger)
	utilTestLevel(t, buf.String(), mlog.Verbose)

	buf.Reset()
	utilFillLogCB(logger)
	utilTestLevel(t, buf.String(), mlog.Verbose)
}

func TestLoggerLevels(t *testing.T) {
	logbookA := nolog.NewLogbook()
	logbookB := nolog.NewLogbook()

	logbook := NewLogbook(logbookA, logbookB)
	logger := logbook.Joiner().Join("test")

	logbook.SetLevel("test", mlog.Fatal)

	if logger.(*teelogger).level != mlog.Fatal {
		t.Error("Expected level to be Fatal")
	}

	logbook.SetLevel(mlog.Default, mlog.Info)

	if logger.(*teelogger).level != mlog.Fatal {
		t.Error("Expected level to be Fatal")
	}

	if logbook.(*teelogbook).defaultlevel != mlog.Info {
		t.Error("Expected level to be Info")
	}
}

func TestLogbookNew(t *testing.T) {
	logbookA := nolog.NewLogbook()
	logbookB := nolog.NewLogbook()

	logbook := NewLogbook(logbookA, logbookB)

	if logbook == nil {
		t.Fatal("Logbook has not been returned")
	}

	joiner := logbook.Joiner()
	if joiner == nil {
		t.Fatal("Joiner has not been returned")
	}

	logger := joiner.Join("test")
	if logger == nil {
		t.Fatal("Logger has not been returned")
	}
}

func TestLogbookLogLevels(t *testing.T) {
	logbookA := nolog.NewLogbook()
	logbookB := nolog.NewLogbook()

	logbook := NewLogbook(logbookA, logbookB)

	if logbook == nil {
		t.Fatal("Logbook has not been returned")
	}

	// Check default levels

	levels := logbook.Levels()

	if len(levels) != 1 {
		t.Fatal("Levels should have only DEFAULT entry instead of", levels)
	}

	if val, ok := levels[mlog.Default]; !ok || val != mlog.Fatal {
		t.Fatal("Level for DEFAULT should be FATAL instead of", levels)
	}

	// Check modified default levels

	if err := logbook.SetLevel(mlog.Default, mlog.Info); err != nil {
		t.Fatal("Error is not expected for SetLevel call")
	}

	levels = logbook.Levels()

	if len(levels) != 1 {
		t.Fatal("Levels should have only DEFAULT entry instead of", levels)
	}

	if val, ok := levels[mlog.Default]; !ok || val != mlog.Info {
		t.Fatal("Level for DEFAULT should be INFO instead of", levels)
	}

	// Check levels for test1 logger

	if err := logbook.SetLevel("test1", mlog.Error); err != nil {
		t.Fatal("Error is not expected for SetLevel call")
	}

	levels = logbook.Levels()

	if len(levels) != 2 {
		t.Fatal("Levels should have DEFAULT & test1 entries instead of", levels)
	}

	if val, ok := levels[mlog.Default]; !ok || val != mlog.Info {
		t.Fatal("Level for DEFAULT should be INFO instead of", levels)
	}

	if val, ok := levels["test1"]; !ok || val != mlog.Error {
		t.Fatal("Level for test1 should be ERROR instead of", levels)
	}

	// Check default for new logger test2

	_ = logbook.Joiner().Join("test2")

	levels = logbook.Levels()

	if len(levels) != 3 {
		t.Fatal("Levels should have DEFAULT & test1 & test2 entries instead of", levels)
	}

	if val, ok := levels[mlog.Default]; !ok || val != mlog.Info {
		t.Fatal("Level for DEFAULT should be INFO instead of", levels)
	}

	if val, ok := levels["test1"]; !ok || val != mlog.Error {
		t.Fatal("Level for test1 should be ERROR instead of", levels)
	}

	if val, ok := levels["test2"]; !ok || val != mlog.Info {
		t.Fatal("Level for test2 should be INFO instead of", levels)
	}

	// Default level change should touch only default dedicated loggers

	if err := logbook.SetLevel(mlog.Default, mlog.Warning); err != nil {
		t.Fatal("Error is not expected for SetLevel call")
	}

	levels = logbook.Levels()

	if len(levels) != 3 {
		t.Fatal("Levels should have DEFAULT & test1 & test2 entries instead of", levels)
	}

	if val, ok := levels[mlog.Default]; !ok || val != mlog.Warning {
		t.Fatal("Level for DEFAULT should be INFO instead of", val)
	}

	if val, ok := levels["test1"]; !ok || val != mlog.Error {
		t.Fatal("Level for test1 should be ERROR instead of", val)
	}

	if val, ok := levels["test2"]; !ok || val != mlog.Warning {
		t.Fatal("Level for test2 should be INFO instead of", val)
	}

	// Loggers with default level should be converted to custom after SetLevel call

	if err := logbook.SetLevel("test2", mlog.Verbose); err != nil {
		t.Fatal("Error is not expected for SetLevel call")
	}

	if err := logbook.SetLevel(mlog.Default, mlog.Fatal); err != nil {
		t.Fatal("Error is not expected for SetLevel call")
	}

	levels = logbook.Levels()

	if len(levels) != 3 {
		t.Fatal("Levels should have DEFAULT & test1 & test2 entries instead of", levels)
	}

	if val, ok := levels[mlog.Default]; !ok || val != mlog.Fatal {
		t.Fatal("Level for DEFAULT should be FATAL instead of", levels)
	}

	if val, ok := levels["test1"]; !ok || val != mlog.Error {
		t.Fatal("Level for test1 should be ERROR instead of", levels)
	}

	if val, ok := levels["test2"]; !ok || val != mlog.Verbose {
		t.Fatal("Level for test2 should be VERBOSE instead of", levels)
	}
}

func TestLogbookJoin(t *testing.T) {
	logbookA := nolog.NewLogbook()
	logbookB := nolog.NewLogbook()

	logbook := NewLogbook(logbookA, logbookB)

	if logbook == nil {
		t.Fatal("Logbook has not been returned")
	}

	joiner := logbook.Joiner()
	if joiner == nil {
		t.Fatal("Joiner has not been returned")
	}

	logger1 := joiner.Join("test")
	if logger1 == nil {
		t.Fatal("Logger has not been returned")
	}

	logger2 := joiner.Join("test")
	if logger2 == nil {
		t.Fatal("Logger has not been returned")
	}

	if logger1 != logger2 {
		t.Fatal("Expected same interface value for same logger")
	}
}
