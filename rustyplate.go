// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// HTTP file system request handler

// This is a derivative of the golang net/http package.
// The main feature as of yet is to allow for custom 404 page handling.


package rustyplate

import (
//	"errors"
	//"fmt"
	//"io"
	//"mime"
//	"mime/multipart"
//	"net/textproto"
//	"net/url"
	//"os"
	"path"
//	"path/filepath"
//	"strconv"
	"strings"
//	"time"
		"net/http"
)


func serveFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string, redirect bool, notFound http.HandlerFunc) {


	f, err := fs.Open(name)
	if err != nil {
		// TODO expose actual error?
		notFound(w, r)
		return
	}
	defer f.Close()

	d, err1 := f.Stat()
	if err1 != nil {
		// TODO expose actual error?
		notFound(w, r)
		return
	}


	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}


type RustyPlate struct {
	root http.FileSystem
	notfound http.HandlerFunc
}


func FileServer(root http.FileSystem)  *RustyPlate {
	return &RustyPlate{root, http.NotFound}
}

func (f *RustyPlate) SetNotFoundFunc(fn http.HandlerFunc ){
	f.notfound = fn
}

func (f *RustyPlate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	serveFile(w, r, f.root, path.Clean(upath), true, f.notfound)
}
