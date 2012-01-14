package main

import (
	"errors"
	"github.com/zeebo/admin"
	"launchpad.net/gobson/bson"
	"time"
)

type Post struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Title   string
	Body    string
	Posted  Time
	Updated Time
}

func (p Post) GetTemplate() string {
	return `<fieldset>
		  <legend>Post</legend>
		  <div class="clearfix{{if .Errors.Title}} error{{end}}">
			<label for="Title">Title</label>
			<div class="input">
			  <input class="xlarge{{if .Errors.Title}} error{{end}}" id="Title" name="Title" size="30" type="text" value="{{.Values.Title}}" />
			  {{if .Errors.Title}}<span class="help-inline">{{.Errors.Title.Error}}</span>{{end}}
			</div>
		  </div><!-- /clearfix -->
		  <div class="clearfix{{if .Errors.Body}} error{{end}}">
			<label for="textarea">Body</label>
			<div class="input">
			  <textarea class="xxlarge{{if .Errors.Body}} error{{end}}" id="Body" name="Body" rows="10">{{.Values.Body}}</textarea>
			  {{if .Errors.Body}}<span class="help-inline">{{.Errors.Body.Error}}</span>{{end}}
			  <span class="help-block">
				Markdown text
			  </span>
			</div>
		  </div><!-- /clearfix -->
		  <div class="actions">
			<input type="submit" class="btn primary" value="Save">
		  </div>
		</fieldset>`
}
func (p *Post) Validate() admin.ValidationErrors {
	errs := make(admin.ValidationErrors)
	if p.Title == "" {
		errs["Title"] = errors.New("Required field.")
	}
	if p.Body == "" {
		errs["Body"] = errors.New("Required field.")
	}
	if p.Posted == 0 {
		p.Posted = Time(time.Now().UnixNano())
	}
	p.Updated = Time(time.Now().UnixNano())
	return errs
}
