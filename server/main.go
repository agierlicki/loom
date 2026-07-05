package main

import (
	"log"
	"net/http"

	"github.com/agierlicki/loom/internal/jobs"
)

func main() {

	jobStore := jobs.NewInMemoryStore()
	jobApi := jobs.NewJobApi(jobStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/jobs", jobApi.JobsHandler)
	mux.HandleFunc("/api/jobs/{id}", jobApi.JobHandler)

	log.Println("Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
