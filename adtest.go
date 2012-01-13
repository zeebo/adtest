package main

import (
	"github.com/zeebo/admin"
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"log"
	"net/http"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(4)
}

func Env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type Post struct {
	ID    bson.ObjectId `bson:"_id"`
	Title string
	Body  string
}

func (p Post) GetTemplate() string {
	return `<fieldset>
          <legend>Post</legend>
          <div class="clearfix">
            <label for="Title">Title</label>
            <div class="input">
              <input class="xlarge" id="Title" name="Title" size="30" type="text" />
            </div>
          </div><!-- /clearfix -->
          <div class="clearfix">
            <label for="textarea">Body</label>
            <div class="input">
              <textarea class="xxlarge" id="Body" name="Body" rows="10"></textarea>
              <span class="help-block">
                Markdown text
              </span>
            </div>
          </div><!-- /clearfix -->
          <div class="actions">
            <input type="submit" class="btn primary" value="Save changes">
          </div>
        </fieldset>`
}
func (p *Post) Validate() admin.ValidationErrors { return nil }

type T2 struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	V  int           `bson:"v"`
}

func (t T2) GetTemplate() string {
	return `<span class="errors">{{.Errors.V.Error}}</span>
	<input type="text" value="{{.Values.V}}" name="V">
	<input type="submit" class="btn primary" value="Submit">`
}
func (t T2) Validate() admin.ValidationErrors { return nil }

func main() {
	session, err := mgo.Mongo(Env("MONGO_URL", "localhost"))
	if err != nil {
		log.Fatal(err)
	}

	a := &admin.Admin{
		Session: session,
		Prefix:  "/admin",
	}
	a.Register(&Post{}, "admin_test.Post", nil)
	a.Register(&T2{}, "admin_test.T2", nil)

	static := http.FileServer(http.Dir(Env("STATIC_DIR", "static")))

	//setup handlers
	http.Handle("/admin/", a)
	http.Handle("/static/", http.StripPrefix("/static/", static))

	if err := http.ListenAndServe(Env("BIND_ADDR", ":11223"), nil); err != nil {
		log.Fatal(err)
	}
}
