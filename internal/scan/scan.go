package scan

import (
    "context"
    "crypto/tls"
    "fmt"
    "io"
    "net"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/PuerkitoBio/goquery"
    "golang.org/x/sync/errgroup"
)

type App struct {
    Name   string
    Port   int
    Scheme string // "http" or "https"
    URL    string
}

const (
    DefaultTimeout = 150 * time.Millisecond
    DefaultWorkers = 2048
)

// Scanner はポート走査に必要な設定とHTTPクライアントを保持します。
type Scanner struct {
    timeout time.Duration
    workers int
    client  *http.Client
}

func New(timeout time.Duration, workers int) *Scanner {
    if workers < 1 {
        workers = 1
    }
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        DialContext:     (&net.Dialer{Timeout: timeout}).DialContext,
    }
    return &Scanner{
        timeout: timeout,
        workers: workers,
        client:  &http.Client{Timeout: timeout * 4, Transport: tr},
    }
}

func NewDefault() *Scanner { return New(DefaultTimeout, DefaultWorkers) }

func PortFromAddr(addr string) (int, bool) {
    parts := strings.Split(addr, ":")
    if len(parts) == 0 {
        return 0, false
    }
    last := parts[len(parts)-1]
    n, err := strconv.Atoi(last)
    if err != nil {
        return 0, false
    }
    return n, true
}

func (s *Scanner) Scan(ctx context.Context, ports []int) []App {
    var (
        g, gctx = errgroup.WithContext(ctx)
        mu       sync.Mutex
        apps     = make([]App, 0, len(ports))
    )
    g.SetLimit(s.workers)

    for _, p := range ports {
        p := p
        g.Go(func() error {
            if app := s.probePort(gctx, p); app != nil {
                mu.Lock()
                apps = append(apps, *app)
                mu.Unlock()
            }
            return nil
        })
    }
    _ = g.Wait()
    return apps
}

func (s *Scanner) probePort(ctx context.Context, port int) *App {
    for _, scheme := range []string{"http", "https"} {
        if name, ok := s.tryFetch(ctx, scheme, port); ok {
            url := fmt.Sprintf("%s://localhost:%d/", scheme, port)
            if name == "" {
                name = url
            }
            return &App{Name: name, Port: port, Scheme: scheme, URL: url}
        }
    }
    return nil
}

// tryFetch: GET / を実行し、応答があればtrue。title取得を試みる。
func (s *Scanner) tryFetch(ctx context.Context, scheme string, port int) (string, bool) {
    url := fmt.Sprintf("%s://localhost:%d/", scheme, port)
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

    resp, err := s.client.Do(req)
    if err != nil {
        return "", false
    }
    defer resp.Body.Close()

    lr := &io.LimitedReader{R: resp.Body, N: 256 * 1024}
    doc, err := goquery.NewDocumentFromReader(lr)
    if err == nil {
        title := strings.TrimSpace(doc.Find("title").First().Text())
        return title, true
    }
    return "", true
}

