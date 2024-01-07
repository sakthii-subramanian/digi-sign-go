package main

import (
    "net/http"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
	"github.com/sakthii-subramanian/digi-sign/signature"
	"github.com/sakthii-subramanian/digi-sign/fileUpload"
	"github.com/sakthii-subramanian/digi-sign/email"
    "github.com/sakthii-subramanian/digi-sign/pdfStore"
)

func main() {
   
    // Set up CORS middleware
    corsMiddleware := handlers.CORS(
        handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Replace with your frontend's URL
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
        handlers.AllowedHeaders([]string{"Content-Type"}),
    )

    // Create a router with the CORS middleware
    r := mux.NewRouter()
    r.Use(corsMiddleware)

    r.HandleFunc("/upload", fileUpload.HandleFileUpload).Methods("POST")
    r.HandleFunc("/signature", signature.HandleSignature).Methods("POST")
    r.HandleFunc("/send-email", email.HandleSendEmail).Methods("POST")
    r.HandleFunc("/pdf/upload",pdfStore.Handle_pdf_db_upload).Methods("POST")
    r.HandleFunc("/pdf/{id:[0-9]+}",pdfStore.Handle_pdf_db_read).Methods("GET")
    r.HandleFunc("/pdf/{id:[0-9]+}",pdfStore.Handle_pdf_db_update).Methods("POST")

    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}
