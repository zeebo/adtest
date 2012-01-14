package main

import (
	"crypto/bcrypt"
	"errors"
	"github.com/zeebo/admin"
	"launchpad.net/gobson/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Username string
	Password string
}

func (u User) GetTemplate() string {
	return `<fieldset>
          <legend>Username</legend>
          <div class="clearfix{{if .Errors.Username}} error{{end}}">
            <label for="Username">Username</label>
            <div class="input">
              <input class="xlarge{{if .Errors.Username}} error{{end}}" id="Username" name="Username" size="30" type="text" value="{{.Values.Username}}" />
              {{if .Errors.Username}}<span class="help-inline">{{.Errors.Username.Error}}</span>{{end}}
            </div>
          </div><!-- /clearfix -->
          <legend>Password</legend>
          <div class="clearfix{{if .Errors.Password}} error{{end}}">
            <label for="Password">Password</label>
            <div class="input">
              <input class="xlarge{{if .Errors.Password}} error{{end}}" id="Password" name="Password" size="30" type="text" value="{{.Values.Password}}" />
              {{if .Errors.Password}}<span class="help-inline">{{.Errors.Password.Error}}</span>{{end}}
            </div>
          </div><!-- /clearfix -->
          <div class="actions">
            <input type="submit" class="btn primary" value="Save">
          </div>
        </fieldset>`
}

func (u *User) Validate() admin.ValidationErrors {
	errs := make(admin.ValidationErrors)
	if u.Username == "" {
		errs["Username"] = errors.New("Required field.")
	}
	if u.Password == "" {
		errs["Password"] = errors.New("Required field.")
	}
	//check for usernames
	n, err := session.DB("blog").C("User").Find(bson.M{"Username": u.Username}).Count()
	if err != nil {
		errs["Username"] = err
	}

	if n > 0 {
		errs["Username"] = errors.New("Username taken")
	}

	if len(errs) == 0 {
		if err := u.HashPassword(); err != nil {
			errs["Password"] = err
		}
	}

	return errs
}

func (u *User) HashPassword() error {
	//a simple heuristic at the moment to determine if it's already hashed
	if u.Password[0] == '$' {
		return nil
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) Compare(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)) == nil
}
