// Copyright (c) 2017 Kuguar <licenses@kuguar.io> Author: Adrian P.K. <apk@kuguar.io>
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package api

import (
	"github.com/adrianpk/fundacja/bootstrap"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Context used for maintaining HTTP Request Context
type Context struct {
	MongoSession *mgo.Session
	User         string
	UserID       bson.ObjectId
}

// Close mgo.Session
func (c *Context) Close() {
	c.MongoSession.Close()
}

// DbCollection returns mgo.collection for the given name
func (c *Context) DbCollection(name string) *mgo.Collection {
	return c.MongoSession.DB(bootstrap.AppConfig.GetDBConnParamenters()["Database"]).C(name)
}

// NewContext creates a new Context object for each HTTP request
func NewContext() *Context {
	session := bootstrap.GetSession().Copy()
	context := &Context{
		MongoSession: session,
	}
	return context
}
