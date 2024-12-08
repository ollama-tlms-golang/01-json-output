package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ollama/ollama/api"
)

var (
	FALSE = false
	TRUE  = true
)

func main() {
	ctx := context.Background()

	var ollamaRawUrl string
	if ollamaRawUrl = os.Getenv("OLLAMA_HOST"); ollamaRawUrl == "" {
		ollamaRawUrl = "http://localhost:11434"
	}

	url, _ := url.Parse(ollamaRawUrl)

	client := api.NewClient(url, http.DefaultClient)

	systemInstructions := `You are a helpful AI assistant. The user will enter the name of an animal.
	The assistant will then return the following information about the animal:
	- the scientific name of the animal (the name of json field is: scientific_name)
	- the main species of the animal  (the name of json field is: main_species)
	- the decimal average length of the animal (the name of json field is: average_length)
	- the decimal average weight of the animal (the name of json field is: average_weight)
	- the decimal average lifespan of the animal (the name of json field is: average_lifespan)
	- the countries where the animal lives into json array of strings (the name of json field is: countries)
	Output the results in JSON format and trim the spaces of the sentence.
	Use the provided context to give the data`

	userContent := "chicken"

	// Prompt construction
	messages := []api.Message{
		{Role: "system", Content: systemInstructions},
		{Role: "user", Content: userContent},
	}

	req := &api.ChatRequest{
		Model:    "granite3-moe:1b",
		Messages: messages,
		Options: map[string]interface{}{
			"temperature":    0.0,
			"repeat_last_n":  2,
		},
		Stream: &FALSE,
		Format: json.RawMessage(`"json"`),
	}

	answer := ""
	err := client.Chat(ctx, req, func(resp api.ChatResponse) error {
		answer = resp.Message.Content
		return nil
	})

	if err != nil {
		log.Fatalln("😡", err)
	}
	fmt.Println(answer)
	fmt.Println()
}
