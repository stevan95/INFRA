package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	server_crt string
	server_key string
}

func NewAPIServer(listenAddr, server_crt, server_key string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		server_crt: server_crt,
		server_key: server_key,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", makeHTTPHandleFunc(s.imageTagController))
	http.ListenAndServeTLS(s.listenAddr, s.server_crt, s.server_key, router)
}

func (s *APIServer) imageTagController(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	if r.ContentLength == 0 {
		return fmt.Errorf("body is empty")
	}

	admissionReview := &AdmissionReview{}
	e := json.NewDecoder(r.Body)
	err := e.Decode(admissionReview)

	if err != nil {
		return fmt.Errorf("cannot covert json")
	}

	// Parse response from go template
	var temp *template.Template
	temp, err = template.ParseFiles("./template/admissionreview.tmpl")
	if err != nil {
		log.Fatalf("Failed to parse template file: %v", err)
	}

	//Get Number of containers running in pod
	for _, container := range admissionReview.Request.Object.Spec.Containers {
		image := strings.Split(container.Image, ":")
		if len(image) > 1 {
			tag := image[1]
			if tag == "latest" {
				admissionRes := AdmissionReviewRes{
					UID:     admissionReview.Request.UID,
					Allowed: false,
					Status: ResponseStatus{
						Code:    403,
						Message: "Using latest tag of the image is not allowed",
					},
				}

				w.Header().Set("Content-Type", "application/json")
				err = temp.Execute(w, admissionRes)
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				admissionRes := AdmissionReviewRes{
					UID:     admissionReview.Request.UID,
					Allowed: true,
					Status: ResponseStatus{
						Code:    200,
						Message: "Image tag is allowed",
					},
				}

				w.Header().Set("Content-Type", "application/json")
				err = temp.Execute(w, admissionRes)
				if err != nil {
					log.Fatalln(err)
				}
			}

		} else {
			admissionRes := AdmissionReviewRes{
				UID:     admissionReview.Request.UID,
				Allowed: false,
				Status: ResponseStatus{
					Code:    403,
					Message: "Using latest tag of the image is not allowed",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			err = temp.Execute(w, admissionRes)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	return nil
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			fmt.Printf("error not cannot convert func to hanle func")
		}
	}
}
