package main

import (
	"github.com/hashlash/descarca/pipe"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Job struct{
	pipe *pipe.Pipe
	consumer <-chan map[string]string
	clientId string
	n int
}

var jobs chan Job

func worker(jobs <-chan Job) {
	for job := range jobs {
		log.Println(job.consumer, job.clientId, job.n)
		err := archiveCompress(job.consumer, job.clientId, job.n)
		if err != nil {
			log.Println("Error while compressing: ", err)
			return
		}
		log.Println("Done archive and compressing")
		job.pipe.Close()
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())
	id := path.Base(r.URL.Path)
	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		log.Println("Can't convert query n to integer: ", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	p, err := pipe.SetupDownloadDonePipeOut(id)
	if err != nil {
		log.Println("Error creating download done pipe: ", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	consumer, err := p.Receive()
	if err != nil {
		log.Println("Cannot retrieve consumer channel: ", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	jobs <- Job{
		pipe: p,
		consumer: consumer,
		clientId: id,
		n: n,
	}
}

func main() {
	jobs = make(chan Job)
	go worker(jobs)

	http.HandleFunc("/", handler)
	log.Println("Serving...")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_HOST"), nil))
}
