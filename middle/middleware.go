package middle

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

func Elasticsearch_connection() *elasticsearch.Client {

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.1.4:9200",
			"elasticuser",
			"user1234",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	return es
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !isLoggedIn(c) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		var loginDetails struct {
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&loginDetails); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body or missing password",
			})
			return
		}

		if !validatePassword(loginDetails.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid password",
			})
			return
		}

		es := Elasticsearch_connection()
		document := struct {
			Userid string `json:"userid"`
		}{
			Userid: c.GetString("userId"),
		}

		data, err := json.Marshal(document)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error marshaling document",
			})
			return
		}

		res, err := es.Index("my_index", bytes.NewReader(data))
		if err != nil || res.IsError() {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to index document in Elasticsearch",
			})
			return
		}
		defer res.Body.Close()

		c.Next()
	}

}

func validatePassword(password string) bool {

	return password == "user1234"

}
