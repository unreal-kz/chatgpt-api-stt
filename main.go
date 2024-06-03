package main

import (
	"context"
	"fmt"
	"log"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	openai "github.com/sashabaranov/go-openai"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func main() {
	fmt.Println(os.Getwd())
	// Load Google Cloud credentials
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "speech-to-text-425309-9686157e4795.json")

	ctx := context.Background()

	// Create a client
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Read the audio file
	data, err := os.ReadFile("path/to/your/audio/file.wav")
	if err != nil {
		log.Fatalf("Failed to read audio file: %v", err)
	}

	// Configure the request
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{
				Content: data,
			},
		},
	}

	// Perform the request
	resp, err := client.Recognize(ctx, req)
	if err != nil {
		log.Fatalf("Failed to recognize: %v", err)
	}

	// Print the results
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("Transcript: %v\n", alt.Transcript)
			// Send the transcript to ChatGPT
			chatGPTResponse := sendToChatGPT(alt.Transcript)
			fmt.Printf("ChatGPT Response: %v\n", chatGPTResponse)
		}
	}
}

func sendToChatGPT(transcript string) string {
	apiKey := "your-openai-api-key"
	client := openai.NewClient(apiKey)

	req := openai.CompletionRequest{
		Model:     "text-davinci-003",
		Prompt:    transcript,
		MaxTokens: 150,
	}

	resp, err := client.CreateCompletion(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get response from ChatGPT: %v", err)
	}

	return resp.Choices[0].Text
}
