package jobs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var tests = []struct {
	name           string
	setupJobID     string
	requestID      string
	expectedStatus int
}{
	{name: "job exists", setupJobID: "foo", requestID: "foo", expectedStatus: http.StatusOK},
	{name: "job not found", setupJobID: "foo", requestID: "bar", expectedStatus: http.StatusNotFound},
}

func TestCreateJobHandler(t *testing.T) {

	store := NewInMemoryStore()
	api := NewJobApi(store)

	body := `{"type": "http_request", "payload": {"url":"http://example.org"}}"`
	post := httptest.NewRequest(http.MethodPost, "/api/jobs", strings.NewReader(body))
	rec := httptest.NewRecorder()

	api.CreateJobHandler(rec, post)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("Invalid response code %d, expected 201", rec.Code)
	}
	var job Job
	if err := json.NewDecoder(res.Body).Decode(&job); err != nil {
		t.Errorf("Response was not a Job struct!")
	}
	if job.Status != JobPending {
		t.Errorf("Resulting job has incorrect status %q", job.Status)
	}
}

func TestGetJobHandler(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := NewInMemoryStore()
			api := NewJobApi(store)

			store.Create(Job{
				ID: test.setupJobID,
			})

			req := httptest.NewRequest(http.MethodGet, "/api/job/{id}", nil)
			req.SetPathValue("id", test.requestID)
			rec := httptest.NewRecorder()

			api.GetJobHandler(rec, req)
			res := rec.Result()

			if res.StatusCode != test.expectedStatus {
				t.Errorf("Response code %d does not match expected code %d", res.StatusCode, test.expectedStatus)
			}
		})
	}
}
