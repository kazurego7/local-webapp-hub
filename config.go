package main

import "flag"

var (
    // 起動時のアドレスのみフラグで指定
    listenAddr = flag.String("addr", ":8787", "listen address, e.g. :8787")
)
