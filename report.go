package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"text/template"
	"time"
)

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
	MultipartForm string    // アップロードファイル情報
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

// Report : Trace.Report を実行することで、トレース情報をレポートとして返却する
func (p *Trace) Report(title interface{}) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("TRACEROUTE").Parse(p.Template(title))
	if err != nil {
		return "", err
	}
	tmpl.Execute(&buf, p)

	return buf.String(), nil
}

// Template : テンプレートメッセージ
func (p *Trace) Template(title interface{}) string {
	message := []string{
		"################################################################################",
		"#",
		"#   CRASH REPORT",
		"#",
		"################################################################################",
		"ErrorTitle:      " + fmt.Sprint(title),
		"ProjectPath:     {{.ProjectDir}}",
		"Binname:         {{.Binname}}({{.Pid}})",
		"Datetime:        " + p.DateTime.Format("2006-01-02 15:04:05"),
		"AccessURL:       {{.AccessURL}}",
		"QueryPath:       {{.Path}}",
		"RemoteAddr:      {{.RemoteAddr}}",
		"UserAgent:       {{.UserAgent}}",
		"Accept:          {{.Accept}}",
		"Accept-Encoding: {{.Encoding}}",
		"Accept-Language: {{.Language}}",
		"{{- if ne .ContentType \"\"}}",
		"Content-Type:    {{.ContentType}}{{end}}",
		"Content-Length:  {{.ContentLength}}",
		"Connection:      {{.Connection}}",
		"RequestMethod:   {{.Method}}",
		"Proto:           {{.ProtoVersion}}",
		"SubmitData: ",
		"{{.Form}}",
		"UploadFiles: ",
		"{{.MultipartForm}}",
		"StackTrace: ",
		"{{range .StackTrace}}{{.}}",
		"{{end}}",
	}
	return strings.Join(message, "\n")
}

// ServeTrace : トレース開始
func ServeTrace(point int, w http.ResponseWriter, r *http.Request) *Trace {
	var p = &Trace{}
	// スタックトレースの取得
	p.StackTrace = []string{}
	for i := point; ; i++ {
		pc, filename, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcname := runtime.FuncForPC(pc).Name()
		p.StackTrace = append(p.StackTrace, fmt.Sprintf("=======>> %d: %s: %s(%d)", i-point, funcname, filename, line))
	}
	// ブラウザ情報
	p.UserAgent = r.UserAgent()
	// リクエストメソッド
	p.Method = r.Method
	// プロジェクトパス
	p.ProjectDir, _ = os.Getwd()
	// 現在時刻
	p.DateTime = time.Now()
	// アクセス元IPアドレス
	p.RemoteAddr = r.RemoteAddr
	// 送信バイト数
	p.ContentLength = r.ContentLength

	// Content-Type
	p.ContentType = r.Header.Get("Content-Type")
	if p.ContentType == "" {
		p.ContentType = w.Header().Get("Content-Type")
	}
	// 送信情報を取得
	if strings.Index(p.ContentType, "multipart/form-data") != -1 {
		r.ParseMultipartForm(32 << 20)
	}
	r.ParseForm()
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		var data = make(map[string]interface{})
		for name, headers := range r.MultipartForm.File {
			var files []map[string]interface{}
			for _, header := range headers {
				var file = map[string]interface{}{
					"filename": header.Filename,
					"headers":  header.Header,
					"filesize": header.Size,
				}
				files = append(files, file)
			}
			data[name] = files
		}
		buf, _ := json.MarshalIndent(data, "", "    ")
		p.MultipartForm = string(buf)
	} else {
		p.MultipartForm = "{}"
	}
	buf, _ := json.MarshalIndent(r.PostForm, "", "    ")
	p.Form = string(buf)
	// ブラウザの言語設定情報
	p.Language = r.Header.Get("Accept-Language")
	// プロトコル名
	p.Protocol = "http"
	if r.TLS != nil {
		p.Protocol = "https"
	}
	p.ProtoVersion = r.Proto
	// アクセスURL
	p.AccessURL = p.Protocol + "://" + r.Host + r.RequestURI
	// ホスト名
	p.Host = r.Host
	// リクエストURI
	p.Path = r.URL.Path
	p.Connection = r.Header.Get("Connection")
	p.Encoding = r.Header.Get("Accept-Encoding")
	p.Accept = r.Header.Get("Accept")
	// プロセスID
	p.Pid = os.Getpid()
	// プログラム名
	p.Binname = os.Args[0]

	return p
}
