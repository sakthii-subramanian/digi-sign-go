// upload.go
package fileUpload

import (
	"net/http"
	"os"
	"io"
	"fmt"

)

func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Print("heyy file upload")
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the uploaded file
	filePath := "./uploads/newfile.pdf" 
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

	w.Write([]byte("File uploaded successfully"))
}

