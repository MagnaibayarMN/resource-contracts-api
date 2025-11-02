package queries

import (
	"encoding/json"
	"fmt"
	"iltodgeree/api/internal/document"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DownloadFile(id string, fileType string, c *gin.Context) {
	docID := uuid.New().String()

	if fileType == "docx" {
		contract, _ := GetContractMaster(id)
		files, err := document.FilePathWalkDir(document.TEMPLATE_PATH)
		document.Check(err)
		document.ProcessSingle(docID, contract, c, files)
	} else {
		contract, _ := GetContract(id)
		var contractJson map[string]interface{}
		_ = json.Unmarshal(*contract.Source, &contractJson)

		fileUrl := contractJson["metadata"].(map[string]interface{})["file_url"].(string)
		re := regexp.MustCompile(`\/([^\/]+\.pdf)$`)
		match := re.FindStringSubmatch(fileUrl)

		if len(match) > 1 {
			fmt.Println("Filename:", match[1])
			path := os.Getenv("STORAGE_PATH") + "/" + contractJson["contract_id"].(string) + "/" + match[1]
			c.File(path)
		} else {
			fmt.Println("No match found")
		}
	}
}
