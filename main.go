package main

import (
	"net/http"
	"log"
	"sync"
	"html/template"
	"path/filepath"
	"flag"
	"chat-go/trace"
	"os"
)

type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.temp1.Execute(w, r)
}

func main() {
	addr := flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	//http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)
	// チャットルームを開始
	go r.run()
	// Webサーバを起動
	log.Println("Webサーバを開始 ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
