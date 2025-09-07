package server

import (
    "net/http"
    "sort"
    "time"

    "github.com/kazurego7/local-webapp-hub/internal/scan"
    "github.com/labstack/echo/v4"
)

// New は Echo を構築し、ルーティングとレンダラーを設定して返します。
func New(scanner *scan.Scanner, listenAddr string) *echo.Echo {
    e := echo.New()
    e.Renderer = newRenderer()

    e.GET("/", func(c echo.Context) error {
        ctx := c.Request().Context()

        // ポート候補は常にフルスキャン（1-65535）
        candidates := make([]int, 0, 65535)
        for p := 1; p <= 65535; p++ {
            candidates = append(candidates, p)
        }

        // 自分自身のポートは除外
        if lp, ok := scan.PortFromAddr(listenAddr); ok {
            filtered := candidates[:0]
            for _, p := range candidates {
                if p != lp {
                    filtered = append(filtered, p)
                }
            }
            candidates = filtered
        }

        start := time.Now()
        apps := scanner.Scan(ctx, candidates)
        sort.Slice(apps, func(i, j int) bool { return apps[i].Port < apps[j].Port })
        dur := time.Since(start)

        data := struct {
            Apps       []scan.App
            Count      int
            DurationMs int64
            Now        time.Time
        }{
            Apps:       apps,
            Count:      len(apps),
            DurationMs: dur.Milliseconds(),
            Now:        time.Now(),
        }

        return c.Render(http.StatusOK, "index.html", data)
    })

    return e
}

