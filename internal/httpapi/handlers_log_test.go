package httpapi

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"buc/internal/app"
	"buc/internal/config"
	"buc/internal/logview"
)

func TestLogFilesIndexAndServe(t *testing.T) {
	manager, err := logview.NewManager(t.TempDir())
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	logview.SetDefault(manager)
	defer logview.SetDefault(nil)

	if _, err := manager.Write([]byte("log line\n")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	router := Router(app.New(&config.AllConfig{}, nil))

	indexReq := httptest.NewRequest(http.MethodGet, "/log/files", nil)
	indexRec := httptest.NewRecorder()
	router.ServeHTTP(indexRec, indexReq)

	if indexRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for files index, got %d", indexRec.Code)
	}
	if !strings.Contains(indexRec.Body.String(), "buc-") {
		t.Fatalf("expected log filename in index, got %q", indexRec.Body.String())
	}

	files, err := manager.Files()
	if err != nil {
		t.Fatalf("Files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 log file, got %d", len(files))
	}

	fileReq := httptest.NewRequest(http.MethodGet, "/log/files/"+files[0].Name, nil)
	fileRec := httptest.NewRecorder()
	router.ServeHTTP(fileRec, fileReq)

	if fileRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for file serve, got %d", fileRec.Code)
	}
	if !strings.Contains(fileRec.Body.String(), "log line") {
		t.Fatalf("expected log contents, got %q", fileRec.Body.String())
	}
}

func TestLiveLogPageAndStream(t *testing.T) {
	manager, err := logview.NewManager(t.TempDir())
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	logview.SetDefault(manager)
	defer logview.SetDefault(nil)

	if _, err := manager.Write([]byte("stream line\n")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	router := Router(app.New(&config.AllConfig{}, nil))

	pageReq := httptest.NewRequest(http.MethodGet, "/log/live", nil)
	pageRec := httptest.NewRecorder()
	router.ServeHTTP(pageRec, pageReq)
	if pageRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for live page, got %d", pageRec.Code)
	}
	if !strings.Contains(pageRec.Body.String(), "EventSource('/api/log/stream')") {
		t.Fatalf("expected EventSource in page, got %q", pageRec.Body.String())
	}

	ctx, cancel := context.WithCancel(context.Background())
	streamReq := httptest.NewRequest(http.MethodGet, "/api/log/stream", nil).WithContext(ctx)
	streamRec := httptest.NewRecorder()
	done := make(chan struct{})
	go func() {
		router.ServeHTTP(streamRec, streamReq)
		close(done)
	}()
	cancel()
	<-done

	if streamRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for stream, got %d", streamRec.Code)
	}
	body, err := io.ReadAll(streamRec.Result().Body)
	if err != nil {
		t.Fatalf("read stream body: %v", err)
	}
	if !strings.Contains(string(body), "stream line") {
		t.Fatalf("expected recent line in stream, got %q", string(body))
	}
}
