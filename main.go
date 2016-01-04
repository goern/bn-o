/*
 Copyright 2016 Christoph Görn

 This file is part of bn-o.

 bn-o is free software: you can redistribute it and/or modify
 it under the terms of the GNU Lesser General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 bn-o is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Lesser General Public License for more details.

 You should have received a copy of the GNU Lesser General Public License
 along with bn-o. If not, see <http://www.gnu.org/licenses/>.
*/

// Package main is the main command line tool for bn-o.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goern/bn-o/Godeps/_workspace/src/github.com/ant0ine/go-json-rest/rest"
	"github.com/goern/bn-o/Godeps/_workspace/src/github.com/coreos/go-semver/semver"
)

type SemVerMiddleware struct {
	MinVersion string
	MaxVersion string
}

func (mw *SemVerMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {

	minVersion, err := semver.NewVersion(mw.MinVersion)
	if err != nil {
		panic(err)
	}

	maxVersion, err := semver.NewVersion(mw.MaxVersion)
	if err != nil {
		panic(err)
	}

	return func(writer rest.ResponseWriter, request *rest.Request) {

		version, err := semver.NewVersion(request.PathParam("version"))
		if err != nil {
			rest.Error(
				writer,
				"Invalid version: "+err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		if version.LessThan(*minVersion) {
			rest.Error(
				writer,
				"Min supported version is "+minVersion.String(),
				http.StatusBadRequest,
			)
			return
		}

		if maxVersion.LessThan(*version) {
			rest.Error(
				writer,
				"Max supported version is "+maxVersion.String(),
				http.StatusBadRequest,
			)
			return
		}

		request.Env["VERSION"] = version
		handler(writer, request)
	}
}

var version string

func main() {
	fmt.Printf("This is bn-o, Version %s\n", version) // FIXME version is not passed along via linker

	svmw := SemVerMiddleware{
		MinVersion: "0.1.0",
		MaxVersion: "0.1.0",
	}

	statusMw := &rest.StatusMiddleware{}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/#version/version", svmw.MiddlewareFunc(
			func(w rest.ResponseWriter, req *rest.Request) {
				// version := req.Env["VERSION"].(*semver.Version)
				w.WriteJson(map[string]string{
					"Name":           "This is bn-o",
					"serviceVersion": version,
				})

			},
		)),
		rest.Get("/#version/status", svmw.MiddlewareFunc(
			func(w rest.ResponseWriter, req *rest.Request) {
				w.WriteJson(statusMw.GetStatus())
			},
		)),
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
