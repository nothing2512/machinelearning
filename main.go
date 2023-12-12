package main

import (
	"fmt"
	"main/controllers"
	"net"
	"net/http"
	"strings"
	"time"
)

const PORT = "3030"

func main() {
	c := controllers.InDb{}

	backends := map[string]http.HandlerFunc{
		"/api/nb/command":  c.NB,
		"/api/knn/command": c.KNN,
		"/api/ann/command": c.ANN,
	}

	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		expires := time.Unix(0, 0)
		w.Header().Set("Expires", expires.Format(http.TimeFormat))
		w.Header().Set("Pragma", "no-cache")
		if handler, ok := backends[r.URL.Path]; ok {
			handler.ServeHTTP(w, r)
			return
		}
		if strings.Contains(r.URL.Path, ".") {
			http.ServeFile(w, r, "static/"+r.URL.Path)
		} else if r.URL.Path[len(r.URL.Path)-1] == '/' {
			http.ServeFile(w, r, "html/"+r.URL.Path+"index.html")
		} else {
			http.ServeFile(w, r, "html/"+r.URL.Path+".html")
		}
	}))

	ln, err := net.Listen("tcp", "0.0.0.0:"+PORT)
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening On 0.0.0.0:" + PORT)
	if err = http.Serve(ln, nil); err != nil {
		panic(err)
	}
}
