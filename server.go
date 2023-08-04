package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	discordwebhook "github.com/bensch777/discord-webhook-golang"
	"github.com/gin-gonic/gin"
)

type Data struct {
	ID         int64  `json:"id"`
	Checkpoint string `json:"checkpoint"`
	Timestamp  int64  `json:"timestamp"`
}

func sendEmbed(link string, embeds discordwebhook.Embed) error {

	hook := discordwebhook.Hook{
		Username:   "TLS23",
		Avatar_url: "https://pss.rs/wp-content/uploads/2023/07/4kolo.jpg",
		Embeds:     []discordwebhook.Embed{embeds},
	}

	payload, err := json.Marshal(hook)
	if err != nil {
		log.Fatal(err)
	}
	err = discordwebhook.ExecuteWebhook(link, payload)
	return err

}

func main() {
	participants := make(map[int64]map[string]int64)
	checkpoints := []string{"KT1", "KT2", "KT3", "KT4", "REV1", "KT5", "REV2", "KT6", "test"}
	var webhookURL = "https://discord.com/api/webhooks/1137042171351007282/gl3MnpHcrm_rN50Y_C-1PyS1BjS7hLkyRb85nN0tgbqPD20TlZxkYAIURgtc87AYJBfj"

	r := gin.Default()

	r.Static("/static", "./static")

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

		embed := discordwebhook.Embed{
			Fields: []discordwebhook.Field{
				discordwebhook.Field{
					Name:   "ID",
					Value:  strconv.FormatInt(data.ID, 10),
					Inline: true,
				},
				discordwebhook.Field{
					Name:   "Kontrolna tačka",
					Value:  data.Checkpoint,
					Inline: true,
				},
			},
			Timestamp: time.Unix(data.Timestamp/1000, 0),
		}

		sendEmbed(webhookURL, embed)

		c.JSON(200, gin.H{"message": fmt.Sprintf("Takmicar %d uspesno prijavljen!", data.ID)})
	})

	r.GET("/api/participants/:id", func(c *gin.Context) {
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

	r.GET("/api/checkpoints/:id", func(c *gin.Context) {
		checkpointID := c.Param("id")

		var data []Data
		for participantID, checkpoints := range participants {
			if timestamp, ok := checkpoints[checkpointID]; ok {
				data = append(data, Data{participantID, checkpointID, timestamp})
			}
		}

		sort.Slice(data, func(i, j int) bool {
			return data[i].Timestamp < data[j].Timestamp
		})

		c.JSON(200, data)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
