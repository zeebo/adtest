package main

import (
	"github.com/zeebo/admin"
	"launchpad.net/mgo"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	session *mgo.Session
)

func init() {
	var err error
	session, err = mgo.Mongo(Env("MONGO_URL", "localhost"))
	if err != nil {
		log.Fatal(err)
	}
}

func Env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type Time int64

func (t Time) String() string {
	return time.Unix(0, int64(t)).String()
}

func main() {
	a := &admin.Admin{
		Session: session,
		Prefix:  "/admin",
	}
	a.Register(&Post{}, "blog.Post", &admin.Options{
		Columns: []string{"Title", "Body", "Posted", "Updated"},
	})
	a.Register(&User{}, "blog.User", &admin.Options{
		Columns: []string{"Username", "Password"},
	})

	static := http.FileServer(http.Dir(Env("STATIC_DIR", "static")))

	//setup handlers
	http.Handle("/admin/", a)
	http.Handle("/static/", http.StripPrefix("/static/", static))
	//http.Handle("/", http.HandleFunc(blog))

	if err := http.ListenAndServe(Env("BIND_ADDR", ":11223"), nil); err != nil {
		log.Fatal(err)
	}
}
