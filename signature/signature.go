// signature.go
package signature

import (
	"net/http"
	"os"
	"io"
	"github.com/sakthii-subramanian/digi-sign/unipdflib"
	"fmt"

)

func HandleSignature(w http.ResponseWriter, r *http.Request) {
	file, fheader, err := r.FormFile("customFile")
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Print("hey handle sign")
	// fmt.Print(r.FormValue("filename"))
	// Save the uploaded file
	filePath := "./uploads/" + fheader.Filename
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error copying the file", http.StatusInternalServerError)
		return
	}

	// Add signature to the PDF
	signaturePath := "./static/signature.png" // Provide the path to your signature image
	if err := unipdflib.SignWithImage(filePath,signaturePath,"/Users/sakthi/digi-sign/static/watermark.png","./uploads/output.pdf"); 
	err != nil {
		http.Error(w, "Error adding signature", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Signature added successfully"))
}
