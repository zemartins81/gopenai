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

type OAIChoices struct {
	Text         string `json:"text"`
	Index        uint8  `json:"index"`
	LogProbs     uint8  `json:"log_probs"`
	FinishReason string `json:"finish_reason"`
}

type OAIResponse struct {
	Id      string       `json:"id"`
	Object  string       `json:"object"`
	Create  uint64       `json:"create"`
	Model   string       `json:"model"`
	Choices []OAIChoices `json:"choices"`
}

type OAIRequest struct {
	Prompt     string `json:"prompt"`
	Max_tokens uint32 `json:"max_tokens"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar arquivo .env")
		return
	}

	fmt.Printf("\x1bc")
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		userInput, _ := reader.ReadString('\n')

		s.Start()
		s.Suffix = "A IA está processando sua requisição"
		requestOpenAi(userInput)
		s.Stop()
	}

}

func requestOpenAi(userInput string) {
	oaiToken := os.Getenv("OPENAI_KEY")
	bearer := "Bearer " + oaiToken
	preamble := `Answer the question in portuguese.`
	uri := "https://api.openai.com/v1/engines/text-davinci-003/completions"

	oaiRquest := OAIRequest{
		Prompt:     fmt.Sprintf("%s %s", preamble, userInput),
		Max_tokens: 1200,
	}

	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(oaiRquest)
	if err != nil {
		log.Fatal(err)
	}

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

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response OAIResponse
	err = json.Unmarshal([]byte(bytes), &response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("")

	fmt.Println(response.Choices[0].Text)
}
