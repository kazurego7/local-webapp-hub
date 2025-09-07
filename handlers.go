package main

import (
    "net/http"
    "sort"
    "time"

    "github.com/labstack/echo/v4"
)

func handleIndex(c echo.Context) error {
    ctx := c.Request().Context()

    // ポート候補は常にフルスキャン（1-65535）
    candidates := make([]int, 0, 65535)
    for p := 1; p <= 65535; p++ {
        candidates = append(candidates, p)
    }

    // 自分自身のポートは除外
    if lp, ok := portFromAddr(*listenAddr); ok {
        filtered := candidates[:0]
        for _, p := range candidates {
            if p != lp {
                filtered = append(filtered, p)
            }
        }
        candidates = filtered
    }

    start := time.Now()
    apps := defaultScanner.Scan(ctx, candidates)
    sort.Slice(apps, func(i, j int) bool { return apps[i].Port < apps[j].Port })
    dur := time.Since(start)

    data := struct {
        Apps       []App
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
}
