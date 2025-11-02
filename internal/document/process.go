// Package document provides utilities for generating and exporting contract documents.
// It supports DOCX generation from templates and TSV/CSV export of search results.
package document

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"iltodgeree/api/internal/correction"
	"iltodgeree/api/internal/sql"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v5"
)

// DOCUMENT_PATH is the base path for storing generated documents.
var DOCUMENT_PATH = ""

// TEMPLATE_PATH is the path to DOCX template files.
var TEMPLATE_PATH = ""

// PUBLIC_URL is the public-facing URL for accessing documents.
var PUBLIC_URL = ""

var GIN_MODE = os.Getenv("GIN_MODE")
var MAIN_DOCUMENT_FILE = "word/document.xml"
var RELEASE = "release"

// ProcessSingle generates a DOCX file from a single contract document.
// It creates a Word document with the contract text and returns it for download.
//
// Parameters:
//   - id: Unique identifier for this export operation
//   - searchResult: The contract document from Elasticsearch
//   - c: Gin context for HTTP response
//   - files: Template files to include in the DOCX
func ProcessSingle(id string, searchResult *elastic.GetResult, c *gin.Context, files []FileBuffer) {
	documentPath := DOCUMENT_PATH + "/result/" + string(id)

	_, err := exec.Command("mkdir", "-p", documentPath).Output()
	Check(err)

	target := "/" + id + ".docx"

	path := documentPath + target
	file, err := os.Create(path)
	Check(err)
	defer file.Close()

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	var dataContents bytes.Buffer
	escapedBuffer := bufio.NewWriter(&dataContents)
	InitializeDocument(escapedBuffer)

	extraLines, err := regexp.Compile("\n\n")
	Check(err)

	var contract map[string]interface{}
	_ = json.Unmarshal(*searchResult.Source, &contract)
	var metadata map[string]interface{} = contract["metadata"].(map[string]interface{})
	contractName := metadata["contract_name"].(string)
	escapedBuffer.WriteString(CreateTitle(contractName))

	singleLined := extraLines.ReplaceAllString(contract["pdf_text_string"].(string), "\n")
	sanitized := strings.Replace(singleLined, "&nbsp;", " ", -1)
	content := strings.Split(sanitized, "\n")
	for _, line := range content {
		escapedBuffer.WriteString(CreateParagraph(XmlEscape(line)))
	}

	escapedBuffer.WriteString(CreateFooter())
	escapedBuffer.Flush()
	Check(err)
	files = append(files, FileBuffer{
		Name: MAIN_DOCUMENT_FILE,
		Data: dataContents.Bytes(),
	})
	for _, file := range files {
		appendZip(zipWriter, file.Name, file.Data)
	}
	err = zipWriter.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Write(zipBuffer.Bytes())

	defer func() {
		err := os.RemoveAll(documentPath)
		if err != nil {
			log.Println("Failed to remove dir: " + documentPath)
			log.Println(err)
		}
	}()
	c.File(path)
}

// Process generates a DOCX file from multiple contract documents.
// It combines multiple contracts into a single Word document with numbered sections.
//
// Parameters:
//   - id: Unique identifier for this export operation
//   - searchResult: Search results containing multiple contracts
//   - c: Gin context for HTTP response
//   - files: Template files to include in the DOCX
func Process(id string, searchResult *elastic.SearchResult, c *gin.Context, files []FileBuffer) {
	documentPath := DOCUMENT_PATH + "/result/" + string(id)

	_, err := exec.Command("mkdir", "-p", documentPath).Output()
	Check(err)

	target := "/" + id + ".docx"

	path := documentPath + target
	file, err := os.Create(path)
	Check(err)
	defer file.Close()

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	var dataContents bytes.Buffer
	escapedBuffer := bufio.NewWriter(&dataContents)
	InitializeDocument(escapedBuffer)

	extraLines, err := regexp.Compile("\n\n")
	Check(err)

	if searchResult.Hits.TotalHits > 0 {
		for index, hit := range searchResult.Hits.Hits {
			var contract map[string]interface{}
			err = json.Unmarshal(*hit.Source, &contract)

			if err != nil {
				continue
			}

			metadata, ok := contract["metadata"].(map[string]interface{})
			if !ok {
				continue
			}
			contractName := metadata["contract_name"].(string)

			escapedBuffer.WriteString(CreateTitle(fmt.Sprintf("%d. %s", index+1, contractName)))

			singleLined := extraLines.ReplaceAllString(contract["pdf_text_string"].(string), "\n")
			sanitized := strings.Replace(singleLined, "&nbsp;", " ", -1)
			content := strings.Split(sanitized, "\n")
			for _, line := range content {
				escapedBuffer.WriteString(CreateParagraph(XmlEscape(line)))
			}
		}
	}
	escapedBuffer.WriteString(CreateFooter())
	escapedBuffer.Flush()
	Check(err)
	files = append(files, FileBuffer{
		Name: MAIN_DOCUMENT_FILE,
		Data: dataContents.Bytes(),
	})
	for _, file := range files {
		appendZip(zipWriter, file.Name, file.Data)
	}
	err = zipWriter.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Write(zipBuffer.Bytes())

	defer func() {
		err := os.RemoveAll(documentPath)
		if err != nil {
			log.Println("Failed to remove dir: " + documentPath)
			log.Println(err)
		}
	}()

	c.File(path)
}

// FileBuffer holds a file's name and binary content in memory.
type FileBuffer struct {
	Name string // File path or name
	Data []byte // File content
}

func appendZip(writer *zip.Writer, name string, data []byte) {
	file, err := writer.Create(strings.Replace(name, fmt.Sprintf("%s/", TEMPLATE_PATH), "", -1))
	if err != nil {
		Check(err)
	}
	_, err = file.Write(data)
	if err != nil {
		Check(err)
	}
}

func convertToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case []string:
		return strings.Join(val, ";")
	case []interface{}:
		// Convert each element to string
		strs := make([]string, len(val))
		for i, item := range val {
			strs[i] = fmt.Sprint(item)
		}
		return strings.Join(strs, ";")
	default:
		return ""
	}
}

// CSV generates a TSV (tab-separated values) file from search results.
// It exports contract metadata, resources, provinces, and annotations.
//
// Parameters:
//   - id: Unique identifier for this export operation
//   - searchResult: Search results to export
//   - c: Gin context for HTTP response
func CSV(id string, searchResult *elastic.SearchResult, c *gin.Context) {

	units, err := sql.GetProvincesAllUnits()
	Check(err)

	documentPath := DOCUMENT_PATH + "/result/" + string(id)
	_, err = exec.Command("mkdir", "-p", documentPath).Output()
	Check(err)

	target := "/" + id + ".tsv"
	filename := documentPath + target

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	// 2. Create a new CSV writer
	writer := csv.NewWriter(file)
	writer.Comma = '\t' // ← Set tab delimiter

	defer writer.Flush()

	// Гэрээний нэр
	// Эрдсийн төрөл
	// Гэрээний төрөл
	// Гэрээ байгуулсан огноо
	// Баримт бичгийн төрөл
	// Гэрээ байгуулсан төрийн байгууллага
	// Компанийн нэр
	// Компанийн РД
	// Компанийн хаяг
	// Байгууллагын харъялал
	// Толгой компани
	// Оролцооны хувь
	// Open Corporates холбоос
	// Тусгай зөвшөрөл эзэмшигч мөн эсэх
	// Тусгай зөвшөөрлийн нэр
	// Тусгай зөвшөөрлийн дугаар
	// Төслийн нэр
	// Төслийн дугаар
	// Эх сурвалж
	// Гэрээний файл
	// Гэрээний тэмдэглэл
	// Аннотацийн төрлүүд
	// Аннотацийн текст
	// OCID

	// 3. Write rows to CSV
	header := []string{
		"#",
		"Гэрээний нэр",
		"Эрдсийн төрөл",
		"Гэрээний төрөл",
		"Гэрээ байгуулсан огноо",
		"Баримт бичгийн төрөл",
		"Аймаг / Сум",
		"Гэрээ байгуулсан төрийн байгууллага",
		"Компанийн нэр",
		"Төслийн нэр",
		"Гэрээний файл",
		"OCID",
		"Аннотацийн текст",
		"Метадата текст",
	}
	var data [][]string

	if searchResult.Hits.TotalHits > 0 {
		for index, hit := range searchResult.Hits.Hits {
			var contract map[string]interface{}
			err = json.Unmarshal(*hit.Source, &contract)

			if err != nil {
				continue
			}

			metadata, ok := contract["metadata"].(map[string]interface{})
			if !ok {
				continue
			}

			provinces := ""
			for _, p := range metadata["provinces"].([]interface{}) {
				u := p.(map[string]interface{})
				province := u["province"].(string)
				district := u["district"].(string)

				pid, _ := strconv.Atoi(province)
				did, _ := strconv.Atoi(district)

				if pid > 0 && did > 0 {
					provinces += units[pid] + " " + units[did] + ";"
				}
			}

			governments := ""
			for _, gov := range metadata["government_entity"].([]interface{}) {
				governments += gov.(map[string]interface{})["entity"].(string) + ";"
			}

			resources := ""
			for _, res := range metadata["resource"].([]interface{}) {
				resources += correction.Resources[res.(string)] + ";"
			}

			metadataString, ok := contract["metadata_string"].(string)
			find := "https://admin.iltodgeree.mn/app"
			// replace := "https://beta-api.iltodgeree.mn/storage"
			replace := PUBLIC_URL + "/storage"

			if ok {
				metadataString = convertToString(strings.Replace(metadataString, find, replace, 2))
			} else {
				metadataString = ""
			}

			data = append(data, []string{
				strconv.Itoa(index+1) + ".",
				metadata["contract_name"].(string),
				resources,
				correction.ContractTypes[convertToString(metadata["contract_type"])],
				convertToString(metadata["signature_date"]),
				correction.DocumentTypes[convertToString(metadata["document_type"])],
				provinces,
				governments,
				convertToString(metadata["company_name"]),
				convertToString(metadata["project_title"]),
				PUBLIC_URL + "/api/contracts/download/" + hit.Id + "/pdf",
				metadata["open_contracting_id"].(string),
				convertToString(contract["annotations_string"]),
				metadataString,
			})
		}
	}

	// Write header
	if err := writer.Write(header); err != nil {
		log.Fatalf("Could not write header: %v", err)
	}

	// Write each row
	for _, row := range data {
		if err := writer.Write(row); err != nil {
			log.Printf("Could not write row: %v", err)
		}
	}

	log.Println("CSV file created successfully")

	defer func() {
		err := os.RemoveAll(filename)
		if err != nil {
			log.Println("Failed to remove dir: " + filename)
			log.Println(err)
		}
	}()

	c.File(filename)
}

// FilePathWalkDir recursively reads all files in a directory into memory.
// Used for loading DOCX template files.
//
// Parameters:
//   - root: Root directory path to walk
//
// Returns:
//   - []FileBuffer: List of files with their content
//   - error: Error if file reading fails
func FilePathWalkDir(root string) ([]FileBuffer, error) {
	var files []FileBuffer

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // Return the error instead of panicking
		}

		if !d.IsDir() {
			dat, err := os.ReadFile(path)
			if err != nil {
				return err // Return the error instead of panicking
			}
			files = append(files, FileBuffer{
				Name: path,
				Data: dat,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

// func processDocument() {
// 	client, err := elastic.NewClient()
// 	if err != nil {
// 		Check(err)
// 	}
// 	files, err := FilePathWalkDir(fmt.Sprintf("./%s", TEMPLATE_PATH))
// 	Check(err)

// 	r := gin.Default()
// 	r.Static("/public", "./result")
// 	r.GET("/get", func(c *gin.Context) {
// 		id := uuid.New().String()

// 		query := c.Query("q")
// 		fmt.Println(query)
// 		searchResult, err := client.Search().
// 			Index("iltodgeree").
// 			Type("master").
// 			Source(query).
// 			Do()
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf(query)
// 		fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
// 		process(id, searchResult, c, files)
// 	})
// 	r.Run()
// }
