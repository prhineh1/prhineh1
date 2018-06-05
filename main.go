package main

import (
	"net/http"

	c "github.com/prhineh1/Panurge/config"
	r "github.com/prhineh1/Panurge/routes"
)

func main() {
	c.SetupEnv()
	http.Handle("/", r.Logger(c.Env, r.Index(c.Env)))
	http.Handle("/login", r.Logger(c.Env, r.Login(c.Env)))
	http.Handle("/register", r.Logger(c.Env, r.Register(c.Env)))
	http.Handle("/logout", r.Logger(c.Env, r.Logout(c.Env)))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	c.Env.Log.Println("Server is starting...")
	http.ListenAndServe(":8080", nil)
}
