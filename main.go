package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, paramName string) (*http.Request, error) {

	paths := []string{"./from_line.jpg"}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	var err error
	req := &http.Request{}

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormFile(paramName, filepath.Base(path))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)

		err = writer.Close()
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest("POST", uri, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		file.Close()
	}

	return req, err
}

func main() {
	// path, _ := os.Getwd()
	// path += "./pepe.png"

	var wg sync.WaitGroup

	wg.Add(10)

	for i := 1; i < 11; i++ {
		go func() {

			request, err := newfileUploadRequest("http://localhost:8080/upload?claimID=123", "photo")
			if err != nil {
				log.Fatal(err)
			}
			client := &http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				log.Fatal(err)
			} else {
				// body := &bytes.Buffer{}
				// _, err := body.ReadFrom(resp.Body)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				resp.Body.Close()
				// fmt.Println(resp.StatusCode)
				// fmt.Println(resp.Header)
				// fmt.Println(body)
			}

			wg.Done()
		}()
	}

	wg.Wait()

}
