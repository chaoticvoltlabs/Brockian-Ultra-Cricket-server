package logview

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	defaultRetentionFiles = 24
	defaultRecentLines    = 200
	logFilePrefix         = "buc-"
	logFileSuffix         = ".log"
)

type FileInfo struct {
	Name    string
	Size    int64
	ModTime time.Time
}

type Manager struct {
	mu             sync.RWMutex
	dir            string
	stdout         io.Writer
	now            func() time.Time
	currentHourKey string
	currentFile    *os.File
	subscribers    map[int]chan string
	nextSubscriber int
	recent         []string
	retentionFiles int
	maxRecentLines int
}

func NewManager(logDir string) (*Manager, error) {
	if strings.TrimSpace(logDir) == "" {
		logDir = "./logs"
	}
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	m := &Manager{
		dir:            logDir,
		stdout:         os.Stdout,
		now:            time.Now,
		subscribers:    map[int]chan string{},
		retentionFiles: defaultRetentionFiles,
		maxRecentLines: defaultRecentLines,
	}

	if err := m.rotateLocked(m.now()); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) Write(p []byte) (int, error) {
	if m == nil {
		return len(p), nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.rotateLocked(m.now()); err != nil {
		return 0, err
	}

	if _, err := m.stdout.Write(p); err != nil {
		return 0, err
	}
	if m.currentFile != nil {
		if _, err := m.currentFile.Write(p); err != nil {
			return 0, err
		}
	}

	m.appendRecentLocked(string(p))
	m.broadcastLocked(string(p))
	return len(p), nil
}

func (m *Manager) Subscribe() (int, <-chan string, []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.nextSubscriber
	m.nextSubscriber++

	ch := make(chan string, 128)
	m.subscribers[id] = ch

	recent := append([]string(nil), m.recent...)
	return id, ch, recent
}

func (m *Manager) Unsubscribe(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch, ok := m.subscribers[id]
	if !ok {
		return
	}
	delete(m.subscribers, id)
	close(ch)
}

func (m *Manager) Files() ([]FileInfo, error) {
	m.mu.RLock()
	dir := m.dir
	m.mu.RUnlock()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !isLogFileName(name) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, FileInfo{
			Name:    name,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name > files[j].Name
	})
	return files, nil
}

func (m *Manager) ReadFile(name string) ([]byte, error) {
	if !isLogFileName(name) {
		return nil, os.ErrNotExist
	}

	m.mu.RLock()
	path := filepath.Join(m.dir, name)
	m.mu.RUnlock()

	return os.ReadFile(path)
}

func (m *Manager) appendRecentLocked(chunk string) {
	if chunk == "" {
		return
	}
	m.recent = append(m.recent, chunk)
	if len(m.recent) > m.maxRecentLines {
		m.recent = append([]string(nil), m.recent[len(m.recent)-m.maxRecentLines:]...)
	}
}

func (m *Manager) broadcastLocked(chunk string) {
	for id, ch := range m.subscribers {
		select {
		case ch <- chunk:
		default:
			delete(m.subscribers, id)
			close(ch)
		}
	}
}

func (m *Manager) rotateLocked(now time.Time) error {
	hourKey := now.Format("2006-01-02-15")
	if hourKey == m.currentHourKey && m.currentFile != nil {
		return nil
	}

	if m.currentFile != nil {
		_ = m.currentFile.Close()
		m.currentFile = nil
	}

	path := filepath.Join(m.dir, logFilePrefix+hourKey+logFileSuffix)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	m.currentFile = f
	m.currentHourKey = hourKey
	if err := m.cleanupOldFilesLocked(); err != nil {
		return err
	}
	return nil
}

func (m *Manager) cleanupOldFilesLocked() error {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if isLogFileName(name) {
			names = append(names, name)
		}
	}

	sort.Strings(names)
	if len(names) <= m.retentionFiles {
		return nil
	}

	for _, name := range names[:len(names)-m.retentionFiles] {
		if err := os.Remove(filepath.Join(m.dir, name)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func isLogFileName(name string) bool {
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return false
	}
	return strings.HasPrefix(name, logFilePrefix) && strings.HasSuffix(name, logFileSuffix)
}
