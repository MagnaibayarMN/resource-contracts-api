// Package main implements the front-end service for the Iltodgeree mining contracts API.
// This service provides RESTful endpoints for searching, retrieving, and managing
// mining contracts, metadata, annotations, and related documents from Elasticsearch.
package main

import (
	"encoding/json"
	"fmt"
	"iltodgeree/api/internal/correction"
	"iltodgeree/api/internal/document"
	"iltodgeree/api/internal/queries"
	"iltodgeree/api/internal/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// CustomRecovery provides a custom panic recovery handler for Gin framework.
// It catches panics during request processing and returns a structured error response
// instead of crashing the server.
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("🔥 Panic caught: %v\n", r)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "something went wrong",
					"panic": fmt.Sprintf("%v", r),
				})
			}
		}()
		c.Next()
	}
}

// main initializes and starts the front-end API service.
// It sets up:
// - Environment variables from .env file
// - PostgreSQL database connection
// - Gin web framework with middleware
// - All HTTP route handlers
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	document.DOCUMENT_PATH = os.Getenv("DOCUMENT_PATH")
	document.TEMPLATE_PATH = os.Getenv("TEMPLATE_PATH")
	document.PUBLIC_URL = os.Getenv("PUBLIC_URL")

	sql.EstablishPgSQL()
	defer sql.Pgsql.Close()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// r.Use(cors.Default())
	r.Use(gin.Logger())
	r.Use(CustomRecovery())
	r.Use(gin.Recovery())

	r.GET("/api/metadata/:id", func(c *gin.Context) {
		id := c.Param("id")

		res, err := queries.GetMetadata(id)
		if *err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/summary", func(c *gin.Context) {
		res, err := queries.Aggregations()

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/summary/year/province/:id", func(c *gin.Context) {
		pId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			panic(err)
		}

		res, err := queries.YearFilterAggregations(pId)

		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/contracts-latest", func(c *gin.Context) {
		res, err := queries.GetLatestContracts(20)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{})
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/provinces/all-units", func(c *gin.Context) {
		provinces, err := sql.GetProvincesAllUnits()

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, provinces)
	})

	r.GET("/api/provinces", func(c *gin.Context) {
		provinces, err := sql.GetProvinces(c.Query("province_id"))

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, provinces)
	})

	r.GET("/api/search", func(c *gin.Context) {
		params := queries.NewSearchParams(
			c.Query("q"),
			c.Query("year"),
			c.Query("contract_type"),
			c.Query("resource"),
			c.Query("company"),
			c.Query("government"),
			c.Query("document_type"),
		)

		if c.Query("province") != "" {
			params.SetProvince(c.Query("province"))
		}

		if c.Query("district") != "" {
			params.SetDistrict(c.Query("district"))
		}

		if c.Query("annotation_category") != "" {
			params.SetAnnotationCategories(c.Query("annotation_category"))
		}

		if c.Query("annotated") != "" {
			ann, err := strconv.ParseBool(c.Query("annotated"))
			if err != nil {
				log.Println("boolean утгыг хөрвүүлж чадсангүй.")
				panic(err)
			}
			params.SetAnnotated(ann)
		}

		params.SetSize(c.Query("size"))
		params.SetFrom(c.Query("from"))

		params.SetSortBy(c.Query("sort_by"))
		params.SetOrder(c.Query("is_asc"))

		res, err := queries.SearchV2(params)
		if *err != nil {
			panic(err)
		}

		if c.Query("download") != "" && c.Query("type") == "docx" {
			files, err := document.FilePathWalkDir(document.TEMPLATE_PATH)
			if err != nil {
				panic(err)
			}
			document.Process(uuid.New().String(), res, c, files)
		} else if c.Query("download") != "" && c.Query("type") == "tsv" {
			document.CSV(uuid.New().String(), res, c)
		} else {
			c.JSON(http.StatusOK, res)
		}
	})

	r.GET("/api/contracts/:id", func(c *gin.Context) {
		id := c.Param("id")

		res, err := queries.GetContract(id)
		if *err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/contracts/:id/text", func(c *gin.Context) {
		id := c.Param("id")

		contract, err := queries.GetContractMaster(id)
		if *err != nil {
			panic(err)
		}

		var contractJson map[string]interface{}
		error := json.Unmarshal(*contract.Source, &contractJson)
		if error != nil {
			panic(error)
		}

		c.JSON(http.StatusOK, map[string]interface{}{"text": contractJson["pdf_text_string"]})
	})

	r.GET("/api/page/:id", func(c *gin.Context) {
		id := c.Param("id")
		locale := c.Query("locale")

		res, e := sql.GetPage(id, locale)

		if e != nil {
			panic(e)
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/law/:id", func(c *gin.Context) {
		id := c.Param("id")
		locale := c.Query("locale")

		res, e := sql.GetLaw(id, locale)
		if e != nil {
			panic(e)
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/contracts/download/:id/:type", func(c *gin.Context) {
		id := c.Param("id")
		fileType := c.Param("type")

		queries.DownloadFile(id, fileType, c)
	})

	r.GET("/api/contracts/:id/annotations", func(c *gin.Context) {
		id := c.Param("id")

		res, e := queries.GetAnnotationByContract(id)

		if *e != nil {
			panic(e)
		}

		c.JSON(http.StatusOK, res)
	})

	r.POST("/api/correction/resources", func(c *gin.Context) {
		index := os.Getenv("ELASTICSEARCH_SECONDARY")
		docType := os.Getenv("ELASTICSEARCH_DOC_MASTER")
		var data []map[string]interface{}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, d := range data {
			key := d["key"]
			value := d["value"]

			correction.ResourcesCorrection(index, docType, key.(string), value.(string))
		}

		c.JSON(http.StatusOK, nil)
	})

	r.POST("/api/correction/contract_types", func(c *gin.Context) {
		index := os.Getenv("ELASTICSEARCH_SECONDARY")
		docType := os.Getenv("ELASTICSEARCH_DOC_MASTER")
		var data []map[string]interface{}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, d := range data {
			key := d["key"]
			value := d["value"]

			correction.ContractTypesCorrection(index, docType, key.(string), value.(string))
		}

		c.JSON(http.StatusOK, nil)
	})

	r.POST("/api/correction/document_types", func(c *gin.Context) {
		index := os.Getenv("ELASTICSEARCH_SECONDARY")
		docType := os.Getenv("ELASTICSEARCH_DOC_MASTER")
		var data []map[string]interface{}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, d := range data {
			key := d["key"]
			value := d["value"]

			correction.DocumentTypesCorrection(index, docType, key.(string), value.(string))
		}

		c.JSON(http.StatusOK, nil)
	})

	r.GET("/storage/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		c.Header("Access-Control-Allow-Origin", os.Getenv("FRONT_END_URL")) // Ensure CORS is set
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		// TODO:
		c.File(os.Getenv("STORAGE_PATH") + "/" + filepath) // Serve file from storage directory
	})

	// r.OPTIONS("/*any", func(c *gin.Context) {
	// 	c.Header("Access-Control-Allow-Origin", "http://localhost:3000") // Ensure CORS is set
	// 	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	// 	c.Status(204) // No Content
	// })

	r.Run()
}
