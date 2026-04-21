package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// ActivityLogger 活动日志记录器
type ActivityLogger struct {
	file       *os.File
	mu         sync.Mutex
	logChan    chan LogEntry
	maxAge     time.Duration // 日志最大保留时间
	cleanupMux sync.Mutex
}

// LogEntry 日志条目
type LogEntry struct {
	Time     string `json:"time"`
	IP       string `json:"ip"`
	License  string `json:"license"`
	DeviceID string `json:"device_id"`
	Success  bool   `json:"success"`
	Reason   string `json:"reason"`
}

// NewActivityLogger 创建日志记录器
func NewActivityLogger(logPath string, maxAge time.Duration) (*ActivityLogger, error) {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	logger := &ActivityLogger{
		file:    f,
		logChan: make(chan LogEntry, 1000), // 缓冲通道
		maxAge:  maxAge,
	}

	// 启动异步写入协程
	go logger.writeLoop()

	// 启动定期清理协程
	go logger.cleanupLoop()

	return logger, nil
}

// Log 发送日志（非阻塞）
func (al *ActivityLogger) Log(entry LogEntry) {
	select {
	case al.logChan <- entry:
	default:
		// 通道满时丢弃日志，不影响主流程
	}
}

// writeLoop 异步写入日志
func (al *ActivityLogger) writeLoop() {
	for entry := range al.logChan {
		data, err := json.Marshal(entry)
		if err != nil {
			continue
		}

		al.mu.Lock()
		_, err = al.file.Write(data)
		if err != nil {
			log.Printf("write activity log: %v", err)
		} else {
			al.file.WriteString("\n")
		}
		al.mu.Unlock()
	}
}

// cleanupLoop 定期清理旧日志
func (al *ActivityLogger) cleanupLoop() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		al.cleanup()
	}
}

// cleanup 清理过期的日志行
func (al *ActivityLogger) cleanup() {
	al.cleanupMux.Lock()
	defer al.cleanupMux.Unlock()

	cutoff := time.Now().Add(-al.maxAge)

	// 读取所有日志行
	al.mu.Lock()
	_, err := al.file.Seek(0, 0)
	if err != nil {
		al.mu.Unlock()
		return
	}

	// 重新打开文件进行读取
	f, err := os.Open(al.file.Name())
	al.mu.Unlock()
	if err != nil {
		return
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	var kept []LogEntry
	for {
		var entry LogEntry
		if err := decoder.Decode(&entry); err != nil {
			break
		}
		t, err := time.Parse(time.RFC3339, entry.Time)
		if err != nil {
			continue
		}
		if !t.Before(cutoff) {
			kept = append(kept, entry)
		}
	}

	// 重写文件
	al.mu.Lock()
	defer al.mu.Unlock()

	al.file.Close()
	f.Close()

	f, err = os.OpenFile(al.file.Name(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("reopen log file: %v", err)
		return
	}

	for _, entry := range kept {
		data, _ := json.Marshal(entry)
		f.Write(data)
		f.WriteString("\n")
	}

	al.file = f
}

// Close 关闭日志器
func (al *ActivityLogger) Close() {
	close(al.logChan)
	al.mu.Lock()
	defer al.mu.Unlock()
	al.file.Close()
}
