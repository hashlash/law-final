package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashlash/descarca/pipe"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./pages/home.html")
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	routingKey := uuid.New()
	println(routingKey.String())

	formKeys := []string{
		"url1", "url2", "url3", "url4", "url5",
		"url6", "url7", "url8", "url9", "url10",
	}

	p, err := pipe.SetupDownloadJobPipe()
	if err != nil {
		log.Println("Cannot create amqp pipe: ", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	defer p.Close()

	u, err := url.Parse(fmt.Sprintf("http://%s/%s", os.Getenv("SERVER3_HOST"), routingKey.String()))
	if err != nil {
		log.Println("Error creating url: ", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	q := u.Query()
	q.Set("n", fmt.Sprint(len(formKeys)))
	u.RawQuery = q.Encode()

	log.Println("Connecting: ", u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("Error while communicating with server 3: ", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("Server 3 internal err")
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully contact server 3 for id %s: %s\n", routingKey, resp.Status)

	for _, key := range formKeys {
		err = p.Send(map[string]string{
			"key": key,
			"uuid": routingKey.String(),
			"url": r.FormValue(key),
		})
		if err != nil {
			log.Println("Failed to send data: ", err)
			continue
		}
	}

	tmpl, err := template.ParseFiles("./pages/progress.html")
	if err != nil {
		log.Println("Failed to parse html file:", err)
		http.Error(w, "Error", http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, map[string]string {
		"stompUrl": os.Getenv("STOMP_URL"),
		"exchangeName": os.Getenv("EXCHANGE_TOPIC"),
		"downloadRoutingKey": pipe.DownloadProgressRoutingKey(routingKey.String()),
		"compressRoutingKey": pipe.CompressProgressRoutingKey(routingKey.String()),
	})
	if err != nil {
		log.Println("Failed to execute template:", err)
		http.Error(w, "Error", http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		switch r.Method {
		case http.MethodGet:
			homePageHandler(w, r)
		case http.MethodPost:
			downloadHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	log.Println("Serving...")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_HOST"), nil))
}
