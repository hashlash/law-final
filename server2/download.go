package main

import (
	"fmt"
	"github.com/hashlash/descarca/pipe"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

type ProgressWriter struct {
	key string
	total, progress int64
	pipe *pipe.Pipe
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.progress += int64(n)
	log.Println(string(p[:10]))
	log.Printf("Read %d bytes for a total of %d from %d total\n", n, pw.progress, pw.total)
	err := pw.pipe.Send(map[string]string{
		"key": pw.key,
		"current": fmt.Sprint(n),
		"progress": fmt.Sprint(pw.progress),
		"total": fmt.Sprint(pw.total),
	})
	if err != nil {
		return 0, err
	}
	return n, nil
}

func download(key, url, id string) (int64, error) {
	log.Println("Downloading: ", url)
	resp, err := http.Get(url)
	if err != nil {
		//log.Println("Error while creating download request: ", err)
		return 0, err
	}
	defer resp.Body.Close()

	dirpath := path.Join(os.Getenv("DOWNLOAD_PATH"), id)
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		err = os.Mkdir(dirpath, os.ModeDir)
		if err != nil {
			//log.Println("Unable to create download folder: ", err)
			return 0, err
		}
	}

	fname := path.Base(resp.Request.URL.String())
	if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "" {
		if _, params, err := mime.ParseMediaType(contentDisposition); err == nil {
			if filename, ok := params["filename"]; ok {
				fname = filename
			}
		}
	}

	file, err := os.Create(path.Join(dirpath, fname))
	if err != nil {
		log.Println("Error while creating download file: ", err)
		return 0, err
	}
	defer file.Close()

	pProgress, err := pipe.SetupDownloadProgressPipeIn(id)
	if err != nil {
		return 0, err
	}
	defer pProgress.Close()

	err = pProgress.Send(map[string]string{
		"key": key,
		"url": url,
		"fname": fname,
	})
	if err != nil {
		return 0, err
	}

	tee := io.TeeReader(resp.Body, &ProgressWriter{key: key, total: resp.ContentLength, pipe: pProgress})
	n, err := io.Copy(file, tee)

	pDone, err := pipe.SetupDownloadDonePipeIn(id)
	if err != nil {
		return 0, err
	}
	defer pDone.Close()

	err = pDone.Send(map[string]string{
		"key": key,
		"fname": fname,
	})
	if err != nil {
		return 0, err
	}

	return n, err
}

