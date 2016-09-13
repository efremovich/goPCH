package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/video-upload", upload)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://192.168.1.14")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Mx-ReqToken,X-Requested-With")
	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
	} else if r.Method == "POST" {
		// body, _ := ioutil.ReadAll(r.Body)
		// fmt.Println(string(body))
		client := GetClient()
		url := "http://192.168.1.14" + r.URL.Path
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header = r.Header
		req.Body = r.Body
		for i, cookie := range r.Cookies() {
			req.Cookies(cookie)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
	fmt.Fprintf(w, "Барев, метод запроса: %s", r.Method)
	fmt.Println(r.Method)
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test1/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
