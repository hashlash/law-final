package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/hashlash/descarca/pipe"
	"io"
	"log"
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

func archiveCompress(consumer <-chan map[string]string, id string, n int) error {
	file, err := os.Create(path.Join(os.Getenv("PUBLIC_FOLDER"), fmt.Sprintf("%s.tar.gz", id)))
	if err != nil {
		return err
	}
	defer file.Close()

	gw, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for i := 0; i < n; i++ {
		log.Printf("Waiting %d-th file\n", i+1)
		data := <- consumer
		log.Println("Compressing: ", data["fname"])
		err := addFile(tw, data["key"], id, data["fname"])
		if err != nil {
			return err
		}
	}
	return nil
}

func addFile(tw *tar.Writer, key, id, fname string) error {
	fpath := path.Join(os.Getenv("DOWNLOAD_PATH"), id, fname)
	file, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	if stat, err := file.Stat(); err == nil{
		hdr, err := tar.FileInfoHeader(stat, stat.Name())
		if err != nil {
			return err
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		p, err := pipe.SetupCompressionProgressQueue(id)
		if err != nil {
			return err
		}

		tee := io.TeeReader(file, &ProgressWriter{key: key, total: stat.Size(), pipe: p})
		if _, err := io.Copy(tw, tee); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
