package main

import (
	"github.com/zeebo/admin"
	"launchpad.net/gobson/bson"
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

func Auth(req *http.Request) (resp admin.AuthResponse) {
	req.ParseForm()
	username, password := req.Form.Get("Username"), req.Form.Get("Password")

	log.Printf("%q %q", username, password)

	var valid User
	if err := session.DB("blog").C("User").Find(bson.M{"username": username}).One(&valid); err != nil {
		resp.Error = err.Error()
		return
	}

	log.Printf("%+v", valid)

	if !valid.ID.Valid() {
		resp.Error = "Invalid username and/or password. Please try again."
		return
	}

	if valid.Compare(password) {
		resp.Passed = true
		resp.Username = valid.Username
		resp.Key = valid.ID
	} else {
		resp.Error = "Invalid username and/or password. Please try again."
	}

	log.Printf("%+v", resp)

	return
}

func main() {
	a := &admin.Admin{
		Session: session,
		Prefix:  "/admin",
		Auth:    admin.AuthFunc(Auth),
	}
	a.Register(&Post{}, "blog.Post", &admin.Options{
		Columns: []string{"Title", "Body", "Posted", "Updated"},
	})
	a.Register(&User{}, "blog.User", &admin.Options{
		Columns: []string{"Username", "Password"},
	})
	a.Init()

	static := http.FileServer(http.Dir(Env("STATIC_DIR", "static")))

	//setup handlers
	http.Handle("/admin/", a)
	http.Handle("/static/", http.StripPrefix("/static/", static))
	//http.Handle("/", http.HandleFunc(blog))

	if err := http.ListenAndServe(Env("BIND_ADDR", ":11223"), nil); err != nil {
		log.Fatal(err)
	}
}
