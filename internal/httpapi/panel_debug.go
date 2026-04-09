package httpapi

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"buc/internal/support"
)

type panelLastSeen struct {
	UpdatedAt  string `json:"updated_at"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	PanelMAC   string `json:"panel_mac"`
	PanelIP    string `json:"panel_ip"`
	RemoteAddr string `json:"remote_addr"`
	UserAgent  string `json:"user_agent"`
}

var (
	panelSeenMu   sync.Mutex
	panelSeenLast panelLastSeen
)

func recordPanelRequest(r *http.Request) panelLastSeen {
	panelMAC := r.Header.Get("X-Panel-MAC")
	if panelMAC == "" {
		panelMAC = r.URL.Query().Get("panel_mac")
	}

	panelIP := r.Header.Get("X-Panel-IP")
	if panelIP == "" {
		panelIP = r.URL.Query().Get("panel_ip")
	}

	info := panelLastSeen{
		UpdatedAt:  time.Now().Format(time.RFC3339),
		Method:     r.Method,
		Path:       r.URL.Path,
		PanelMAC:   panelMAC,
		PanelIP:    panelIP,
		RemoteAddr: r.RemoteAddr,
		UserAgent:  r.UserAgent(),
	}

	panelSeenMu.Lock()
	panelSeenLast = info
	panelSeenMu.Unlock()

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		log.Printf("panel request method=%s path=%s panel_mac=%s panel_ip=%s remote_ip=%s",
			info.Method, info.Path, info.PanelMAC, info.PanelIP, host)
	} else {
		log.Printf("panel request method=%s path=%s panel_mac=%s panel_ip=%s remote=%s",
			info.Method, info.Path, info.PanelMAC, info.PanelIP, info.RemoteAddr)
	}

	return info
}

func PanelDebugLastSeenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panelSeenMu.Lock()
		info := panelSeenLast
		panelSeenMu.Unlock()

		support.JSON(w, http.StatusOK, info)
	}
}
