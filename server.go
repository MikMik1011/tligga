package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Data struct {
	ID         int64  `json:"id"`
	Checkpoint string `json:"checkpoint"`
	Timestamp  int64  `json:"timestamp"`
}

func main() {
	participants := make(map[int64]map[string]int64)
	checkpoints := []string{"KT1"}

	r := gin.Default()

	r.POST("/api/participant", func(c *gin.Context) {
		var data Data
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(400, gin.H{"error": "Failed to parse JSON"})
			return
		}

		if _, ok := participants[data.ID]; !ok {
			participants[data.ID] = make(map[string]int64)
		}

		participants[data.ID][data.Checkpoint] = data.Timestamp

		c.JSON(200, gin.H{"message": "Data added to the nested map."})
	})

	r.GET("/api/participant/:id", func(c *gin.Context) {
		participantID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid participant ID"})
			return
		}

		if data, ok := participants[participantID]; ok {
			c.JSON(200, data)
		} else {
			c.JSON(404, gin.H{"error": "Participant not found"})
		}
	})

	r.GET("/api/checkpoints", func(c *gin.Context) {
		c.JSON(200, checkpoints)
	})

	r.GET("/api/checkpoint/:id", func(c *gin.Context) {
		checkpointID := c.Param("id")

		var data []Data
		for participantID, checkpoints := range participants {
			if timestamp, ok := checkpoints[checkpointID]; ok {
				data = append(data, Data{participantID, checkpointID, timestamp})
			}
		}

		c.JSON(200, data)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
