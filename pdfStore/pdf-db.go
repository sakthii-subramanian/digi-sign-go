package pdfStore

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	dbUser     = "root"
	dbPassword = "Lucky@2002"
	dbName     = "digisign"
	uploadPath = "./uploads" // Change this to the path where you want to store uploaded PDFs
)

func handle_pdf_db() {
	// Connect to the MySQL database.
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a table to store PDF files.
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS pdf_files (
			id INT AUTO_INCREMENT PRIMARY KEY,
			pdf_content LONGBLOB
		);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}

	// Create the uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		panic(err)
	}

}


func Handle_pdf_db_upload(w http.ResponseWriter, r *http.Request) {
	// Handle PDF upload
	
	handle_pdf_db()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	fmt.Print("db connection success \n")

		file, handler, err := r.FormFile("pdf")
		if err != nil {
			http.Error(w, "Error parsing PDF file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save the PDF file to the server
		pdfPath := filepath.Join(uploadPath, handler.Filename)
		out, err := os.Create(pdfPath)
		if err != nil {
			http.Error(w, "Error saving PDF file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Error saving PDF file", http.StatusInternalServerError)
			return
		}

		// Store the PDF file path in the database
		insertQuery := "INSERT INTO pdf_files (pdf_content) VALUES (?)"
		result, err := db.Exec(insertQuery, pdfPath)
		if err != nil {
			http.Error(w, "Error storing PDF file in the database", http.StatusInternalServerError)
			return
		}

		// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
    w.Write([]byte(fmt.Sprintf("%d", lastInsertID)))
	}


	// Read PDF by ID
func Handle_pdf_db_read(w http.ResponseWriter, r *http.Request) {
	handle_pdf_db()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if err != nil {
		panic(err)
	}
	defer db.Close()
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid PDF ID", http.StatusBadRequest)
			return
		}

		// Retrieve the PDF file path from the database
		var pdfPath string
		err = db.QueryRow("SELECT pdf_content FROM pdf_files WHERE id = ?", id).Scan(&pdfPath)
		if err != nil {
			http.Error(w, "PDF not found", http.StatusNotFound)
			return
		}

		// Send the PDF file as the response
		http.ServeFile(w, r, pdfPath)
	}

	// Update PDF by ID
	func Handle_pdf_db_update(w http.ResponseWriter, r *http.Request) {

		handle_pdf_db()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if err != nil {
		panic(err)
	}
	defer db.Close()
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid PDF ID", http.StatusBadRequest)
			return
		}
		
		// Retrieve the PDF file path from the database
		var pdfPath string
		err = db.QueryRow("SELECT pdf_content FROM pdf_files WHERE id = ?", id).Scan(&pdfPath)
		if err != nil {
			http.Error(w, "PDF not found", http.StatusNotFound)
			return
		}
	
		// // Remove the existing PDF file
		// if err := os.Remove(pdfPath); err != nil {
		// 	http.Error(w, "Error removing existing PDF file", http.StatusInternalServerError)
		// 	return
		// }

		// Save the new PDF file to the server
		file, _, err := r.FormFile("pdf")
		if err != nil {
			http.Error(w, "Error parsing new PDF file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		
		// Save the PDF file to the server
		newPdfPath := filepath.Join(uploadPath, fmt.Sprintf("new_%s", filepath.Base(pdfPath)))
		out, err := os.Create(newPdfPath)
		if err != nil {
			http.Error(w, "Error saving new PDF file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Error saving new PDF file", http.StatusInternalServerError)
			return
		}

		// Update the PDF file path in the database
		updateQuery := "UPDATE pdf_files SET pdf_content = ? WHERE id = ?"
		_, err = db.Exec(updateQuery, newPdfPath, id)
		if err != nil {
			http.Error(w, "Error updating PDF file in the database", http.StatusInternalServerError)
			return
		}

		
		w.WriteHeader(http.StatusCreated)
    	w.Write([]byte(fmt.Sprintf("%d", id)))
	}

