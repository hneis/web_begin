package main

import (
	"encoding/json"
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
	YANDEX_OAUTH = os.Getenv("YANDEX_OAUTH")
	header       = map[string][]string{
		"Accept":        {"application/json"},
		"Content-Type":  {"application/json"},
		"Authorization": {"OAuth " + YANDEX_OAUTH},
	}
)

type jsonResponse struct {
	Href      string
	Method    string
	Templated bool
}

type Download struct {
	Destination string
	client      http.Client
}

/*
@filename to download
*/
func (d Download) GetFileFromYandexDisk(filename string) error {
	url, err := url.Parse(fmt.Sprintf("https://cloud-api.yandex.net/v1/disk/resources/download?path=/%s", filename))
	if err != nil {
		log.Println(err)
		return err
	}

	jsResp, err := d.queryToDownlod(url)
	if err != nil {
		return err
	}

	url, err = url.Parse(jsResp.Href)
	if err != nil {
		log.Println(err)
		return err
	}

	err = d.download(jsResp.Method, url)
	if err != nil {
		return err
	}

	return nil
}

func (d Download) queryToDownlod(url *url.URL) (jr *jsonResponse, err error) {
	r, err := d.client.Do(&http.Request{
		Method: "GET",
		URL:    url,
		Header: header,
	})

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
	r, err := d.client.Do(&http.Request{
		Method: m,
		URL:    url,
		Header: header,
	})

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

func main() {
	if YANDEX_OAUTH == "" {
		fmt.Println("Variable YANDEX_OAUTH is not pressed")
		return
	}
	downloader := Download{
		Destination: "/tmp/",
	}
	err := downloader.GetFileFromYandexDisk("test.jpg")
	if err != nil {
		log.Println(err)
	}

	fmt.Println("File test.jpg saved in " + downloader.Destination)
}
