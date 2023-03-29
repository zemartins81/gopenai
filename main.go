package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	
	"github.com/briandowns/spinner"
	"github.com/joho/godotenv"
)

// Define a struct for the Choices object in the response
type OAIChoices struct {
	Text         string `json:"text"`
	Index        uint8  `json:"index"`
	LogProbs     uint8  `json:"log_probs"`
	FinishReason string `json:"finish_reason"`
}

// Define a struct for the entire response object
type OAIResponse struct {
	Id      string       `json:"id"`
	Object  string       `json:"object"`
	Create  uint64       `json:"create"`
	Model   string       `json:"model"`
	Choices []OAIChoices `json:"choices"`
}

// Define a struct for the request object
type OAIRequest struct {
	Prompt     string `json:"prompt"`
	Max_tokens uint32 `json:"max_tokens"`
}

func main() {
	// Load the OpenAI API key from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar arquivo .env")
		return
	}
	
	fmt.Printf("\x1bc")
	
	// Initialize the spinner animation
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	
	// Initialize reader to read user input
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print(">")
		userInput, _ := reader.ReadString('\n')
		
		// Start the spinner and show a loading message
		s.Start()
		s.Suffix = "A IA está processando sua requisição"
		
		// Call the function to make the OpenAI API request
		requestOpenAi(userInput)
		
		// Stop the spinner animation after the API request completes
		s.Stop()
	}
}

// This function makes a request to the OpenAI API to generate a response
// to the user's input, and prints the response to the console.
func requestOpenAi(userInput string) {
	
	// Get OpenAI API key from .env file
	oaiToken := os.Getenv("OPENAI_KEY")
	
	// Add authentication header
	bearer := "Bearer " + oaiToken
	
	// Set up API request
	preamble := `Answer the question in portuguese.`
	uri := "https://api.openai.com/v1/engines/text-davinci-003/completions"
	oaiRquest := OAIRequest{
		Prompt:     fmt.Sprintf("%s %s", preamble, userInput),
		Max_tokens: 1200,
	}
	
	// Encode the request data as JSON
	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(oaiRquest)
	if err != nil {
		log.Fatal(err)
	}
	
	// Send the request to the OpenAI API
	req, err := http.NewRequest(http.MethodPost, uri, &payload)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	
	// Read the response JSON and decode it
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var response OAIResponse
	err = json.Unmarshal([]byte(bytes), &response)
	if err != nil {
		log.Fatal(err)
	}
	
	// Print the generated response to the console
	fmt.Print("")
	fmt.Println(response.Choices[0].Text)
}
