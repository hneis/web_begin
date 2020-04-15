package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
)

var (
	header = map[string][]string{
		"Accept":       {"application/json"},
		"Content-Type": {"application/json"},
	}
)

type jsonResponse struct {
	Href      string
	Method    string
	Templated bool
}

type Download struct {
	Destination string
	Url         *url.URL
	client      http.Client
}

/*
@filename to download
*/
func (d Download) GetFileFromYandexDisk() error {
	jsResp, err := d.queryToDownlod(d.Url)
	if err != nil {
		return err
	}

	durl, err := url.Parse(jsResp.Href)
	if err != nil {
		log.Println(err)
		return err
	}

	err = d.download(jsResp.Method, durl)
	if err != nil {
		return err
	}

	return nil
}

func (d Download) queryToDownlod(url *url.URL) (jr *jsonResponse, err error) {
	r, err := d.client.Get(url.String())

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&jr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return jr, nil
}

func (d Download) download(m string, url *url.URL) error {
	r, err := d.client.Get(url.String())
	if err != nil {
		log.Println(err)
		return err
	}
	defer r.Body.Close()

	cd := r.Header.Get("Content-Disposition")
	if cd == "" {
		return fmt.Errorf("Content-Disposition is empty")
	}

	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		log.Println(err)
		return err
	}

	filename := params["filename"]
	if filename == "" {
		return fmt.Errorf("filename is empty")
	}

	o, err := os.Create(path.Join(d.Destination, filename))
	if err != nil {
		return err
	}

	_, err = io.Copy(o, r.Body)
	if err != nil {
		return err
	}

	return nil
}

var urlFlag string
var dstDir string

func init() {
	flag.StringVar(&urlFlag, "url", "", "yandex disk shared file link")
	flag.StringVar(&dstDir, "dstdir", "/tmp/", "destination directory")
}

func main() {
	flag.Parse()
	if urlFlag == "" {
		flag.Usage()
	}
	url, err := url.Parse(urlFlag)
	if err != nil {
		log.Println(err)
	}

	downloader := Download{
		Destination: dstDir,
		Url:         url,
	}

	err = downloader.GetFileFromYandexDisk()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("File test.jpg saved in " + downloader.Destination)
}
