package main

import (
	"encoding/json"
	"io/ioutil"
	"k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
)

func main() {
	certFile := "certs/webhook-server-tls.crt"
	keyFile := "certs/webhook-server-tls.key"

	http.HandleFunc("/", admit)
	if err := http.ListenAndServeTLS(":8080", certFile, keyFile, nil); err != nil {
		log.Fatal(err)
	}
}

// admit handles Kubernetes admission requests
func admit(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// Read the request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		setErr(w, "failed to read body")
		return
	}

	// Unmarshal the k8s admission review
	admissionReview := v1.AdmissionReview{}
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		setErr(w, "failed to unmarshall body")
		return
	}

	// This admission controller only follow pods created by k8s admin
	if admissionReview.Request.Kind.Kind != "Pod" || admissionReview.Request.Resource.Resource != "pods" &&
		admissionReview.Request.Operation != "CREATE" || admissionReview.Request.UserInfo.Username != "kubernetes-admin" {
		return
	}

	// Block creation of pods named nginx
	admissionReview.Response = &v1.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: admissionReview.Request.Name != "nginx",
	}

	if !admissionReview.Response.Allowed {
		admissionReview.Response.Result = &metav1.Status{
			Message: "nginx is not allowed today",
			Code:    http.StatusForbidden,
		}
	}

	// Encode and send the admission review response
	resp, err := json.Marshal(admissionReview)
	if err != nil {
		setErr(w, err.Error())
		return
	}

	w.Write(resp)
}

// setErr logs and returns the given error to the given writer
func setErr(w http.ResponseWriter, err string) {
	log.Println(err)
	http.Error(w, err, http.StatusBadRequest)
}
