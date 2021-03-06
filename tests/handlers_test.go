package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prhineh1/panurge/models"

	"github.com/prhineh1/panurge/routes"
	"github.com/stretchr/testify/assert"
)

var rec *httptest.ResponseRecorder
var req *http.Request
var doc bytes.Buffer

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	// Authenticated user
	c := &http.Cookie{
		Name:  "session",
		Value: "active",
	}
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	req.AddCookie(c)
	routes.Index(env).ServeHTTP(rec, req)
	env.Tpl.ExecuteTemplate(&doc, "index.html", models.Session{true})
	assert.Equal(doc.String(), rec.Body.String())
	assert.Equal(200, rec.Result().StatusCode)

	// Non-authenticated user
	doc.Reset()
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	routes.Index(env).ServeHTTP(rec, req)
	env.Tpl.ExecuteTemplate(&doc, "index.html", nil)
	assert.Equal(doc.String(), rec.Body.String())
	assert.Equal(200, rec.Result().StatusCode)
}

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	var c *http.Cookie

	// Authenticated user
	c = &http.Cookie{
		Name:  "session",
		Value: "success",
	}
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/Register", nil)
	req.AddCookie(c)
	routes.Register(env).ServeHTTP(rec, req)
	assert.Equal(303, rec.Result().StatusCode)

	// GET
	c = &http.Cookie{
		Name:  "session",
		Value: "failure",
	}
	doc.Reset()
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/Register", nil)
	req.AddCookie(c)
	env.Tpl.ExecuteTemplate(&doc, "register.html", nil)
	routes.Register(env).ServeHTTP(rec, req)
	assert.Equal(doc.String(), rec.Body.String())
	assert.Equal(200, rec.Result().StatusCode)

	// POST: invalid username
	doc.Reset()
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Register?username=thisnamewillbewaytoolong&password=Abc123?!&email=register@example.com", nil)
	req.AddCookie(c)
	routes.Register(env).ServeHTTP(rec, req)
	env.Tpl.ExecuteTemplate(&doc, "register.html", routes.Message{"invalid username"})
	assert.Equal(doc.String(), rec.Body.String())
	assert.Equal(200, rec.Result().StatusCode)

	// POST: invalid password
	doc.Reset()
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Register?username=registerPost&password=ab()23&email=register@example.com", nil)
	req.AddCookie(c)
	routes.Register(env).ServeHTTP(rec, req)
	env.Tpl.ExecuteTemplate(&doc, "register.html", routes.Message{"invalid password"})
	assert.Equal(doc.String(), rec.Body.String())
	assert.Equal(200, rec.Result().StatusCode)

	// POST: invalid email
	doc.Reset()
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Register?username=registerPost&password=Abc123?!&email=email.example.com", nil)
	req.AddCookie(c)
	routes.Register(env).ServeHTTP(rec, req)
	env.Tpl.ExecuteTemplate(&doc, "register.html", routes.Message{"invalid email"})
	assert.Equal(doc.String(), rec.Body.String())
	assert.Equal(200, rec.Result().StatusCode)

	// POST: success
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Register?username=registerPost&password=Abc123?!&email=register@example.com", nil)
	req.AddCookie(c)
	routes.Register(env).ServeHTTP(rec, req)
	cs := rec.Result().Cookies()
	assert.Equal("idForCreateSession", cs[0].Value)
	assert.Equal(303, rec.Result().StatusCode)

	// POST: duplicate username
	doc.Reset()
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Register?username=duplicateName&password=Abc123?!&email=register@example.com", nil)
	req.AddCookie(c)
	routes.Register(env).ServeHTTP(rec, req)
	env.Tpl.ExecuteTemplate(&doc, "register.html", routes.Message{"duplicateName is taken"})
	assert.Equal(doc.String(), rec.Body.String())

	// POST: duplicate email
	doc.Reset()
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Register?username=duplicateEmail&password=Abc123?!&email=testuser1@mailinator.com", nil)
	req.AddCookie(c)
	routes.Register(env).ServeHTTP(rec, req)
	env.Tpl.ExecuteTemplate(&doc, "register.html", routes.Message{"testuser1@mailinator.com is already in use"})
	assert.Equal(doc.String(), rec.Body.String())

	// POST: 500 errors
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Register?username=createUser500&password=Abc123?!&email=register@example.com", nil)
	req.AddCookie(c)
	routes.Register(env).ServeHTTP(rec, req)
	assert.Equal(500, rec.Result().StatusCode)

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Register?username=createSess500&password=Abc123?!&email=register@example.com", nil)
	req.AddCookie(c)
	routes.Register(env).ServeHTTP(rec, req)
	assert.Equal(500, rec.Result().StatusCode)
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	var c *http.Cookie

	// Authorized user
	c = &http.Cookie{
		Name:  "session",
		Value: "success",
	}
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/Login", nil)
	req.AddCookie(c)
	routes.Login(env).ServeHTTP(rec, req)
	assert.Equal(303, rec.Result().StatusCode)

	// GET
	c = &http.Cookie{
		Name:  "session",
		Value: "failure",
	}
	doc.Reset()
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/Login", nil)
	req.AddCookie(c)
	routes.Login(env).ServeHTTP(rec, req)
	env.Tpl.ExecuteTemplate(&doc, "login.html", nil)
	assert.Equal(doc.String(), rec.Body.String())
	assert.Equal(200, rec.Result().StatusCode)

	// POST confirm cookie exists
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Login?username=loginPost&password=Abc123?!&persist=false", nil)
	req.AddCookie(c)
	routes.Login(env).ServeHTTP(rec, req)
	cs := rec.Result().Cookies()
	assert.Equal("idForCreateSession", cs[0].Value)
	assert.Equal(303, rec.Result().StatusCode)

	// POST incorrect password
	doc.Reset()
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Login?username=loginPost&password=incorrect&persist=false", nil)
	req.AddCookie(c)
	routes.Login(env).ServeHTTP(rec, req)
	env.Tpl.ExecuteTemplate(&doc, "login.html", routes.Message{"Incorrect Password."})
	assert.Equal(doc.String(), rec.Body.String())
	assert.Equal(200, rec.Result().StatusCode)

	// POST incorrect username
	doc.Reset()
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Login?username=noUser&Abc123?!&persist=false", nil)
	req.AddCookie(c)
	routes.Login(env).ServeHTTP(rec, req)
	env.Tpl.ExecuteTemplate(&doc, "login.html", routes.Message{"Username is incorrect."})
	assert.Equal(doc.String(), rec.Body.String())
	assert.Equal(200, rec.Result().StatusCode)

	// 500 error on CreateSession
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/Login?username=createSess500&password=Abc123?!&persist=false", nil)
	req.AddCookie(c)
	routes.Login(env).ServeHTTP(rec, req)
	assert.Equal(500, rec.Result().StatusCode)
}

func TestLogout(t *testing.T) {
	assert := assert.New(t)
	var c *http.Cookie

	// session should have ended
	rec = httptest.NewRecorder()
	rec.Header().Set("Authorization", "registered")
	req, _ = http.NewRequest("GET", "/Logout", nil)
	c = &http.Cookie{
		Name:  "session",
		Value: "end",
	}
	req.AddCookie(c)
	routes.Logout(env).ServeHTTP(rec, req)
	cs := rec.Result().Cookies()
	assert.Equal(303, rec.Result().StatusCode)
	assert.Equal(-1, cs[0].MaxAge)
}

func TestGame(t *testing.T) {
	assert := assert.New(t)
	var c *http.Cookie

	// unauthorized user
	c = &http.Cookie{
		Name:  "session",
		Value: "failure",
	}
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game", nil)
	req.AddCookie(c)
	routes.Game(env).ServeHTTP(rec, req)
	assert.Equal(303, rec.Result().StatusCode)

	// successfuly GET
	c = &http.Cookie{
		Name:  "session",
		Value: "success",
	}
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/game", nil)
	req.AddCookie(c)
	routes.Game(env).ServeHTTP(rec, req)
	doc.Reset()
	env.Tpl.ExecuteTemplate(&doc, "game.html", nil)
	assert.Equal(doc.String(), rec.Body.String())
	assert.Equal(200, rec.Result().StatusCode)
}
