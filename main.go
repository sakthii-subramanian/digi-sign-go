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
    // r := mux.NewRouter()

    // // Enable CORS for all routes.
    // // Replace "*" with the specific origin(s) you want to allow.
    // corsMiddleware := func(next http.Handler) http.Handler {
    //     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //         w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
    //         w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    //         w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    //         if r.Method == "OPTIONS" {
    //             w.WriteHeader(http.StatusNoContent)
    //             return
    //         }

    //         next.ServeHTTP(w, r)
    //     })
    // }

    // // Apply the CORS middleware to all routes.
    // r.Use(corsMiddleware)
   
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
