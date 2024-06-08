package main

import "net/http"

type MethodHandler struct {
	handler http.Handler
	method  string
}

func (mh *MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != mh.method {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	mh.handler.ServeHTTP(w, r)
}

func EnsureMethod(method string, f func(w http.ResponseWriter, r *http.Request)) *MethodHandler {
	return &MethodHandler{
		handler: http.HandlerFunc(f),
		method:  method,
	}
}
