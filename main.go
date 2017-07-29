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
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.temp1.Execute(w, data)
}

const (
	facebookClientId = "566880443702209"
	githubClientId   = "e966db44e197cb80fee7"
	googleClientId   = "613930549056-ma8v8hqdmd6hfbb4kq3diql6ca2tk5kp.apps.googleusercontent.com"
)

func main() {
	addr := flag.String("addr", ":8080", "アプリケーションのアドレス")
	securityKey := flag.String("key", "", "セキュリティキー")
	facebookClientSecret := flag.String("facebook", "", "Facebookクライアントシークレット")
	githubClientSecret := flag.String("github", "", "GitHubクライアントシークレット")
	googleClientSecret := flag.String("google", "", "Googleクラインアントシークレット")
	flag.Parse() // フラグを解釈

	gomniauth.SetSecurityKey(*securityKey)
	gomniauth.WithProviders(
		facebook.New(facebookClientId, *facebookClientSecret, "http://localhost:8080/auth/callback/facebook"),
		github.New(githubClientId, *githubClientSecret, "http://localhost:8080/auth/callback/github"),
		google.New(googleClientId, *googleClientSecret, "http://localhost:8080/auth/callback/google"),
	)

	// チャットルームを作成
	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	// http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	// チャットルームを開始
	go r.run()

	// Webサーバを起動
	log.Println("Webサーバを開始 ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
