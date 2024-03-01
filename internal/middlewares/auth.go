package middlewares

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"todo/internal/session"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" || r.URL.Path == "/register" {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			log.Println("session_id cookie not found")

			templ, err := template.ParseFiles("./frontend/login.html")
			if err != nil {
				log.Printf("не удалось распарсить шаблон: %s\n", err.Error())
			}

			if err = templ.Execute(w, nil); err != nil {
				log.Printf("не удалось выполнить шаблон: %s\n", err.Error())
			}

			return
		}
		sessionID := cookie.Value
		session, ok := session.Sessions[sessionID]
		if !ok {
			log.Println(err)

			templ, err := template.ParseFiles("./frontend/login.html")
			if err != nil {
				log.Printf("не удалось распарсить шаблон: %s\n", err.Error())
			}

			if err = templ.Execute(w, nil); err != nil {
				log.Printf("не удалось выполнить шаблон: %s\n", err.Error())
			}

			return
		}
		userLogin := session.Login
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userLogin", userLogin)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
