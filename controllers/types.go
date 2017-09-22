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

package controllers

const (
	successAlert = "success"
	infoAlert    = "info"
	warningAlert = "warning"
	errorAlert   = "error"
)

// Page - Container for model to render in the template and an alert message.
type Page struct {
	Model interface{}
	Alert *PageAlert
}

func makePage(model interface{}, alert *PageAlert) *Page {
	return &Page{model, alert}
}

// PageAlert - Categorizable alert messages.
type PageAlert struct {
	message string
	kind    string
}

func makePageAlert(message, kind string) *PageAlert {
	var k string
	if k != successAlert && k != infoAlert && k != warningAlert && k != errorAlert {
		k = infoAlert
	}
	return &PageAlert{message, k}
}

func makePageAlertFromError(err error, kind string) *PageAlert {
	return makePageAlert(err.Error(), kind)
}

// ParentChild - Container for Parent and child objects.
type ParentChild struct {
	Parent interface{}
	Child  interface{}
}

func makeParentChild(parent, child interface{}) *ParentChild {
	return &ParentChild{parent, child}
}