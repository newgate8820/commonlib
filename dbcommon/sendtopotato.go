package dbcommon

import (
	"bytes"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap/zapcore"
)

/*************************
* author: Dev0026
* createTime: 18-12-6
* updateTime: 18-12-6
* description:
*************************/

var allSendToPotato = []*SendToPotato{}

type SendToPotato struct {
	currentTime       int64
	maxPerSec         int32
	sendTime          int32
	sendMsg           func(msg string)
	prefix            string
	enableLeve        zapcore.Level
	startWrite        bool
	enableLevel       zapcore.Level
	tempMessages      []string
	tempMessagesMutex sync.RWMutex
}

func NewSendToPotato(maxPerSec int32, fn func(string)) *SendToPotato {
	hostName, _ := os.Hostname()
	stp := &SendToPotato{
		currentTime: time.Now().Unix(),
		maxPerSec:   maxPerSec,
		sendMsg:     fn,
		enableLevel: zapcore.ErrorLevel,
		prefix:      hostName + " \n",
	}
	allSendToPotato = append(allSendToPotato, stp)
	return stp
}

func NewSendToPotatoWithLevel(maxPerSec int32, enableLeve zapcore.Level, fn func(string)) *SendToPotato {
	hostName, _ := os.Hostname()
	stp := &SendToPotato{
		currentTime: time.Now().Unix(),
		maxPerSec:   maxPerSec,
		sendMsg:     fn,
		enableLevel: enableLeve,
		prefix:      hostName + " \n",
	}
	allSendToPotato = append(allSendToPotato, stp)
	return stp
}

func StartWrite(ok bool) {
	for _, one := range allSendToPotato {
		one.startWrite = ok
		for _, msg := range one.tempMessages {
			one.sendMsg(msg)
		}
		one.tempMessages = one.tempMessages[:0]
	}
}

func (s *SendToPotato) EnableLevel(level zapcore.Level) {
	s.enableLevel = level
}

func (s *SendToPotato) Write(b []byte) (n int, er error) {
	var l zapcore.Level

	if len(b) >= 16 && nil == l.UnmarshalText(bytes.Split(b[10:16], []byte{'"'})[0]) && !s.enableLevel.Enabled(l) {
		return len(b), nil
	}
	t := time.Now().Unix()
	if t != s.currentTime {
		atomic.StoreInt64(&s.currentTime, t)
		atomic.StoreInt32(&s.sendTime, 0)
	}
	msg := string(b)
	if val := atomic.AddInt32(&s.sendTime, 1); val == s.maxPerSec {
		msg += "\n\n too many call per sec"
	} else if val > s.maxPerSec {
		return
	}
	if s.startWrite {
		go s.sendMsg(msg)
	} else {
		s.tempMessagesMutex.Lock()
		s.tempMessages = append(s.tempMessages, msg)
		if len(s.tempMessages) > int(s.maxPerSec) {
			s.tempMessages = s.tempMessages[:s.maxPerSec] // 最多缓存20条
		}
		s.tempMessagesMutex.Unlock()
	}
	return len(b), nil
}

func (s *SendToPotato) Sync() error {
	return nil
}
