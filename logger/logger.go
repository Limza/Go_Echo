// Package logger 는 기존 log 패키지를 확장하여 사용하기 위해 만들었다
// write 작업을 할때 mutex 동기화를하기 때문에 고루틴에 안전하다
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// RotateOp 는 log를 저장하는 회전 주기 설정
type RotateOp int

// RotateOp 의 옵션들
const (
	Second RotateOp = iota
	Minute
	Hour
	Day
)

// Logger 는 logger 패키지를 나타내는 구조체
type Logger struct {
	Option   RotateOp
	FileDir  string // 파일이 저장되는 위치
	FileName string // 파일 이름

	file       *os.File
	mu         sync.Mutex
	beforeTime time.Time
}

// NewLogger 는 새로운 Logger를 반환한다
func NewLogger(option RotateOp, fileDir, fileName string) (*Logger, error) {
	logger := &Logger{
		Option:     option,
		FileDir:    fileDir,
		FileName:   fileName,
		beforeTime: time.Now(),
	}

	if err := logger.rotate(); err != nil {
		return nil, err
	}

	log.SetOutput(io.MultiWriter(logger, os.Stdout))

	return logger, nil
}

// Debug 는 디버깅에 관련된 로그를 기록한다
func (l *Logger) Debug(v ...interface{}) {
	log.Println("[DEBUG]", fmt.Sprint(v...))
}

// DebugF 는 디버깅에 관련된 로그를 서식에 맞춰 기록한다
func (l *Logger) DebugF(format string, v ...interface{}) {
	log.Printf("[DEBUG]"+format, v...)
}

// Info 는 정보에 관련된 로그를 기록한다
func (l *Logger) Info(v ...interface{}) {
	log.Println("[INFO]", fmt.Sprint(v...))
}

// InfoF 는 정보에 관련된 로그를 서식에 맞춰 기록한다
func (l *Logger) InfoF(format string, v ...interface{}) {
	log.Printf("[INFO]"+format, v...)
}

// Error 는 에러 로그 기록
func (l *Logger) Error(v ...interface{}) {
	log.Println("[ERROR]", fmt.Sprint(v...))
}

// ErrorF 는 에러 로그를 서식에 맞춰 기록
func (l *Logger) ErrorF(format string, v ...interface{}) {
	log.Printf("[ERROR]"+format, v...)
}

// Fatal 는 중대한 에러를 기록하고 프로그램을 종료시킨다
func (l *Logger) Fatal(v ...interface{}) {
	log.Fatal("[FATAL]", fmt.Sprint(v...))
}

// FatalF 는 중대한 에러를 서식에 맞춰 기록하고 프로그램을 종료시킨다
func (l *Logger) FatalF(format string, v ...interface{}) {
	log.Fatalf("[FATAL]"+format, v...)
}

// Write 는 log 패키지의 SetOutput에 Logger를 등록하기 위해 구현된 인터페이스
// Write를 구현하여 Logger를 동작하게 했다
func (l *Logger) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	diff, err := l.isDiff()
	if err != nil {
		return 0, err
	}

	if l.file == nil || diff {
		l.beforeTime = time.Now()
		l.rotate()
	}

	return l.file.Write(p)
}

func (l *Logger) isDiff() (bool, error) {
	now := time.Now()
	switch l.Option {
	case Second:
		return l.beforeTime.Second() != now.Second(), nil
	case Minute:
		return l.beforeTime.Minute() != now.Minute(), nil
	case Hour:
		return l.beforeTime.Hour() != now.Hour(), nil
	case Day:
		return l.beforeTime.Day() != now.Day(), nil
	default:
		return false, fmt.Errorf("Logger invalid Option")
	}
}

func (l *Logger) rotate() error {
	if err := l.close(); err != nil {
		return err
	}
	if err := l.open(); err != nil {
		return err
	}
	return nil
}

func (l *Logger) open() error {
	if l.FileDir == "" {
		return fmt.Errorf("Logger open func requires FilePDir")
	}
	if l.FileName == "" {
		return fmt.Errorf("Logger open function requires FileName")
	}

	if !strings.HasSuffix(l.FileDir, "/") {
		l.FileDir += "/"
	}
	if _, err := os.Stat(l.FileDir); err != nil {
		if err := os.Mkdir(l.FileDir, os.ModePerm); err != nil {
			return fmt.Errorf("can't make directories for new logfile: %s", err)
		}
	}

	now := time.Now()
	fileFmt := ""
	switch l.Option {
	case Second:
		fileFmt = fmt.Sprintf("_%d-%02d-%02d_%02d;%02d;%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	case Minute:
		fileFmt = fmt.Sprintf("_%d-%02d-%02d_%02d;%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	case Hour:
		fileFmt = fmt.Sprintf("_%d-%02d-%02d_%02d", now.Year(), now.Month(), now.Day(), now.Hour())
	case Day:
		fileFmt = fmt.Sprintf("_%d-%02d-%02d", now.Year(), now.Month(), now.Day())
	}
	fileName := l.FileDir + l.FileName + fileFmt + ".log"
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}

	l.file = f
	return nil
}

func (l *Logger) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}
