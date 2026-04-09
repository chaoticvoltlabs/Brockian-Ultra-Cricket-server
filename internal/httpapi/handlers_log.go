package httpapi

import (
	"fmt"
	"html"
	"net/http"
	"path"
	"strings"
	"time"

	"buc/internal/logview"
)

func LiveLogPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>BUC Live Log</title>
  <style>
    :root { color-scheme: dark; }
    body { margin: 0; background: #101214; color: #d9dde3; font: 14px/1.45 ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; }
    header { padding: 12px 16px; border-bottom: 1px solid #2a2f36; position: sticky; top: 0; background: #101214; }
    a { color: #8dc2ff; text-decoration: none; }
    a:hover { text-decoration: underline; }
    #log { white-space: pre-wrap; padding: 16px; margin: 0; min-height: calc(100vh - 60px); }
  </style>
</head>
<body>
  <header>
    <strong>BUC Live Log</strong>
    <span style="margin-left:12px"><a href="/log/files">Recent log files</a></span>
  </header>
  <pre id="log"></pre>
  <script>
    const el = document.getElementById('log');
    const stream = new EventSource('/api/log/stream');
    function append(text) {
      el.textContent += text;
      window.scrollTo({ top: document.body.scrollHeight, behavior: 'instant' });
    }
    stream.onmessage = (event) => append(event.data + "\n");
    stream.addEventListener('chunk', (event) => append(event.data + "\n"));
    stream.onerror = () => append("[stream disconnected]\n");
  </script>
</body>
</html>`)
	}
}

func LogStreamHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		manager := logview.Default()
		if manager == nil {
			http.Error(w, "log stream unavailable", http.StatusServiceUnavailable)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		id, ch, recent := manager.Subscribe()
		defer manager.Unsubscribe(id)

		for _, chunk := range recent {
			writeSSE(w, "chunk", chunk)
		}
		flusher.Flush()

		heartbeat := time.NewTicker(20 * time.Second)
		defer heartbeat.Stop()

		for {
			select {
			case <-r.Context().Done():
				return
			case chunk, ok := <-ch:
				if !ok {
					return
				}
				writeSSE(w, "chunk", chunk)
				flusher.Flush()
			case <-heartbeat.C:
				fmt.Fprintf(w, ": keepalive\n\n")
				flusher.Flush()
			}
		}
	}
}

func LogFilesIndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		manager := logview.Default()
		if manager == nil {
			http.Error(w, "log files unavailable", http.StatusServiceUnavailable)
			return
		}

		files, err := manager.Files()
		if err != nil {
			http.Error(w, "list log files failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>BUC Log Files</title>
  <style>
    :root { color-scheme: dark; }
    body { margin: 0; background: #101214; color: #d9dde3; font: 14px/1.45 ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; padding: 16px; }
    a { color: #8dc2ff; text-decoration: none; }
    a:hover { text-decoration: underline; }
    table { border-collapse: collapse; width: 100%%; max-width: 960px; }
    th, td { text-align: left; padding: 8px 10px; border-bottom: 1px solid #2a2f36; }
  </style>
</head>
<body>
  <p><a href="/log/live">Live log</a></p>
  <table>
    <thead>
      <tr><th>File</th><th>Modified</th><th>Size</th></tr>
    </thead>
    <tbody>`)
		for _, file := range files {
			fmt.Fprintf(w,
				`<tr><td><a href="/log/files/%s">%s</a></td><td>%s</td><td>%d</td></tr>`,
				html.EscapeString(file.Name),
				html.EscapeString(file.Name),
				html.EscapeString(file.ModTime.Format(time.RFC3339)),
				file.Size,
			)
		}
		fmt.Fprintf(w, `</tbody></table></body></html>`)
	}
}

func LogFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		manager := logview.Default()
		if manager == nil {
			http.Error(w, "log files unavailable", http.StatusServiceUnavailable)
			return
		}

		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/log/files/")
		if name == "." || name == "" || name == "log/files" {
			http.Redirect(w, r, "/log/files", http.StatusFound)
			return
		}

		data, err := manager.ReadFile(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write(data)
	}
}

func writeSSE(w http.ResponseWriter, event string, chunk string) {
	if event != "" {
		fmt.Fprintf(w, "event: %s\n", event)
	}
	for _, line := range strings.Split(strings.ReplaceAll(chunk, "\r\n", "\n"), "\n") {
		fmt.Fprintf(w, "data: %s\n", line)
	}
	fmt.Fprint(w, "\n")
}
