package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type PostParams struct {
	fileType string
	fileSize int
	fileName string
	err      error
	tmpFile  string
	response []byte
	filePath string
	nodePath string
}

func (p *PostParams) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postParams := new(PostParams)
	postParams.nodePath = "/storage/node" + r.Host[len(r.Host)-1:]
	handler(w, r, postParams)
}

func main() {
	// http.HandleFunc("/", handler)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe(":9091", &PostParams{}))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe(":9092", &PostParams{}))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe(":9093", &PostParams{}))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe(":9094", &PostParams{}))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe(":9095", &PostParams{}))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe(":9096", &PostParams{}))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe(":9097", &PostParams{}))
		wg.Done()
	}()
	wg.Wait()
	//http.HandleFunc("/video-upload", upload)
	//log.Fatal(http.ListenAndServe(":9090-9095", nil))

}

func handler(w http.ResponseWriter, r *http.Request, postParams *PostParams) {
	w.Header().Set("Access-Control-Allow-Origin", r.Host)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Mx-ReqToken,X-Requested-With")
	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
	} else if r.Method == "POST" {
		saveTempDataOnNode(r, postParams)

		getPathFromService(r, postParams)

		if _, err := w.Write(postParams.response); err != nil {
			fmt.Printf("WriteFile failed %q\n", err)
		}
		if err := copyFile(postParams.tmpFile, postParams.nodePath+postParams.filePath); err != nil {
			fmt.Printf("CopyFile failed %q\n", err)
		}
	}

}

func copyFile(src, dst string) (err error) {

	fName := filepath.Dir(dst)
	if err := os.MkdirAll(fName, 0644); err != nil {
		fmt.Println(err)
	}
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// Write data to dst
	if err = ioutil.WriteFile(dst, data, 0644); err != nil {
		return err
	}
	return nil
}

func saveTempDataOnNode(msg *http.Request, postParams *PostParams) {

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				postParams.err = err

			}
			// if err != nil {
			// 	log.Println(err)
			// }
			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				log.Println(err)
			}
			// fmt.Printf("Part %v: %v\n", p.Header.Get("Content-Disposition"), slurp)
			if strings.HasPrefix(p.Header.Get("Content-Disposition"), "form-data; name=\"UserReportForm[video]\"; filename=") {
				header := strings.Split(p.Header.Get("Content-Disposition"), "; ")
				postParams.fileName = strings.Trim(strings.Trim(header[2], "filename="), "\"")
				postParams.fileType = p.Header.Get("Content-Type")

				upload(postParams.fileName, slurp, postParams)
			}

		}
	}
}

func getPathFromService(r *http.Request, postParams *PostParams) {
	client := GetClient()

	params := url.Values{}
	params.Set("name", postParams.fileName)
	params.Add("size", string(postParams.fileSize))
	params.Add("type", postParams.fileType)

	// urlStr := "http://192.168.1.14" + r.URL.Path + "?name=" + postParams.filename + "&size=" + string(postParams.filesize) + "&type=" + postParams.filetype
	urlStr := "http://192.168.1.14" + r.URL.Path + "?" + params.Encode()
	req, err := http.NewRequest("POST", urlStr, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header = r.Header

	rawstr := strings.Split(r.URL.RawQuery, "=")
	cookie := http.Cookie{Name: rawstr[0], Value: rawstr[1]}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	qunt := body[:4]
	i, err := strconv.Atoi(string(qunt))
	if err != nil {
		fmt.Println(err)
	}
	postParams.filePath = string(body[4 : i+4])
	postParams.response = body[i+4:]
	// fmt.Printf("%v\n %v\n %v\n", string(qunt), postParams.pathFile, postParams.response)

}

// upload logic
func upload(filename string, data []byte, postParams *PostParams) {
	tmpfile, err := ioutil.TempFile(os.TempDir(), filename)
	size, err := tmpfile.Write(data)
	if err != nil {
		log.Println(err)
	}

	postParams.fileSize = size
	postParams.tmpFile = tmpfile.Name()

}
