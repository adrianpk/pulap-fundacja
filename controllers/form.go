package controllers

import (
	"net/http"

	"github.com/adrianpk/pulap/logger"
	"github.com/gorilla/schema"
)

var (
	formTrue = "true"
)

func getFormDecoder(ignoreUnknownKeys bool) *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(ignoreUnknownKeys)
	return decoder
}

func inputIsTrue(r *http.Request, field string) bool {
	logger.Debugf("Remember is %s", r.PostFormValue("remember"))
	return r.PostFormValue(field) == formTrue
}
