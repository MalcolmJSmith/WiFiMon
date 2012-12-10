package main

import (
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
	"runtime"
	"strings"
	"strconv"
)
const tp = `.\html\`

var (
	chStatus chan WiFiStatus	
	chMon chan WiFiStatus
	chAbout  chan AboutInfo	
	chExit  chan int8
	
)
var templates = template.Must(template.ParseFiles(tp+"wifimon.html", tp+"graph.html", tp+"graph_ie.html",tp+"license.html",tp+"exit.html",tp+"style.css"))

func serveTemplate(w http.ResponseWriter,  templateName string, data interface{}) {

	w.Header().Set("Content-Type", mime.TypeByExtension(templateName[strings.LastIndex(templateName, "."):]))
	w.Header().Set("Cache-Control", "max-age=0")

	err := templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newClient(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		serveTemplate(w, "wifimon.html", "")
	} else {
		http.NotFound(w, r)
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[9:10] == "Y" {
		serveTemplate(w, "exit.html", "")
		chExit <- 0
	} else {
		if r.URL.Path[8:9] == "Y" {
			serveTemplate(w, "license.html",<-chAbout)
		} else {	
			if r.URL.Path[10:11] == "Y" {
				serveTemplate(w, "graph.html", <-chStatus)
			} else  {
				serveTemplate(w, "graph_ie.html", <-chStatus)
			}
		}
	}
}

func style(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, "style.css", "")
}

func main() {

	var(  port uint64
			e error
		)
		
	switch l := len(os.Args); {
	case l < 2:
		port = 80
	case l == 2:
		port,e = strconv.ParseUint(os.Args[1],0,16)
		if e != nil {
			port = 0
		}
	default:
}
	if port == 0 {
		fmt.Println( "wifimon [port]"  )
		os.Exit(87)
	}
	
	fmt.Println("Listening on port", port)	
		
	if runtime.NumCPU() > 1 {
		runtime.GOMAXPROCS(2)
	}
	
	chStatus = make(chan WiFiStatus)
	chMon = make(chan WiFiStatus)
	chAbout = make(chan AboutInfo)
	chExit = make(chan int8)
	
	go Monitor()
	go MessageLoop()
	
	http.HandleFunc("/", newClient)
	http.HandleFunc("/update/", update)
	http.HandleFunc("/style/", style)
	http.ListenAndServe(string(strconv.AppendUint([]byte(":"),port,10)), nil)
}
