package report

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"text/template"
)

func Test__PANIC_DUMP_SUCCESS(t *testing.T) {
	// https サーバを起動
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		// トレース開始
		case "/":
			trace := ServeTrace(0, w, r)
			rep, err := trace.Report("TRACE")
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(rep)
		// カスタムトレース開始
		case "/custom":
			var trace = &MyTrace{
				Trace: ServeTrace(1, w, r),
			}
			rep, err := trace.Report("PANIC")
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(rep)
		}
	}))

	defer ts.Close()

	// リクエストを送信する
	client := ts.Client()
	client.Get(ts.URL)
	client.Get(ts.URL + "/custom")
}

// PANICダンプ用にカスタマイズ
type MyTrace struct {
	*Trace
	ErrorTitle string
}

func (t *MyTrace) Report(title interface{}) (string, error) {
	t.ErrorTitle = fmt.Sprint(title)
	var buf bytes.Buffer
	tmpl, err := template.New("PANICDUMP").Parse(t.Template())
	if err != nil {
		return "", err
	}
	tmpl.Execute(&buf, t)

	return buf.String(), nil
}

func (t *MyTrace) Template() string {
	strs := strings.Split(t.Trace.Template(t.ErrorTitle), "\n")

	res := strings.Join(strs[5:], "\n")
	mes := "################################################################################\n"
	mes += "#\n"
	mes += "#   PANIC DUMP\n"
	mes += "#\n"
	mes += "################################################################################\n"
	res = mes + res
	return res
}
