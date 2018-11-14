トレースライブラリ
===
受け付けたリクエストの詳細(トレース)を取得するライブラリです。

サンプル
---

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/ochipin/report"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                p := report.ServeTrace(1, w, r)
                rep, _ := p.Report(err)
                fmt.Println(rep)
            }
        }()
        w.WriteHeader(200)
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(""))

        panic("ERROR")
    })  
    http.ListenAndServe(":8080", nil)
}
```

`StackTrace`関数は、次の実装になっている。

```go
func ServeTrace(point int, w http.ResponseWriter, r *http.Request) *Trace {...}
```

* 第1引数 - int  
トレース開始位置。基本的に 0 でOK
* 第2引数 - http.ResponseWriter  
* 第3引数 - *http.Request

トレース情報は、`ServeTrace`関数が返却する`Trace`構造体に格納されている。

```go
// Trace : Trace情報を管理する構造体
type Trace struct {
	StackTrace    []string  // スタックトレース
	UserAgent     string    // ブラウザ情報
	Method        string    // リクエストメソッド
	ProjectDir    string    // プロジェクトパス
	DateTime      time.Time // 現在時刻
	RemoteAddr    string    // アクセス元IPアドレス
	ContentLength int64     // 送信バイト数
	AccessURL     string    // アクセスURL
	Form          string    // 送信データ情報
	ContentType   string    // Content-Type
	Language      string    // Accept-Language
	Protocol      string    // http or https などのプロトコル名
	Host          string    // ホスト名
	Path          string    // クエリパス
	ProtoVersion  string    // プロトコルのバージョン
	Pid           int       // プロセスID
	Connection    string    // Connection
	Accept        string    // Accept
	Encoding      string    // Accept-Encoding
	Binname       string    // プログラム名
	ErrorMessage  string    // エラーメッセージ
}
```

`Trace.Report`関数を実行することで、トレース情報を取得することが可能。

```go
// トレース情報を取得
trace := report.StackTrace(1, w, r)
rep, _ := trace.Report("ERROR_TITLE")
// トレース情報を出力
fmt.Println(rep)
// ################################################################################
// #
// #   CRASH REPORT
// #
// ################################################################################
// ErrorTitle:      ERROR_TITLE
// ProjectPath:     /home/user/app
// Binname:         bin/www(2091)
// Datetime:        2018-11-09 17:31:37
// AccessURL:       http://localhost:8080/path/to/url?name=other
// QueryPath:       /path/to/url
// RemoteAddr:      127.0.0.1:39862
// UserAgent:       Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:63.0) Gecko/20100101 Firefox/63.0
// Accept:          text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
// Accept-Encoding: gzip, deflate
// Accept-Language: ja,en-US;q=0.7,en;q=0.3
// Content-Type:    text/html
// Content-Length:  0
// Connection:      keep-alive
// RequestMethod:   POST
// Proto:           HTTP/1.1
// SubmitData: 
// {
//     "OK": [
//         true
//     ]
// }
// StackTrace: 
// =======>> 0: main.main.func1.1: /home/user/app/app.go(14)
// =======>> 1: runtime.call32: /home/user/.go/src/go1.11/src/runtime/asm_amd64.s(522)
// =======>> 2: runtime.gopanic: /home/user/.go/src/go1.11/src/runtime/panic.go(513)
// =======>> 3: main.main.func1: /home/user/app/app.go(23)
// =======>> 4: net/http.HandlerFunc.ServeHTTP: /home/user/.go/src/go1.11/src/net/http/server.go(1964)
// =======>> 5: net/http.(*ServeMux).ServeHTTP: /home/user/.go/src/go1.11/src/net/http/server.go(2361)
// =======>> 6: net/http.serverHandler.ServeHTTP: /home/user/.go/src/go1.11/src/net/http/server.go(2741)
// =======>> 7: net/http.(*conn).serve: /home/user/.go/src/go1.11/src/net/http/server.go(1847)
// =======>> 8: runtime.goexit: /home/user/.go/src/go1.11/src/runtime/asm_amd64.s(1333)
```
