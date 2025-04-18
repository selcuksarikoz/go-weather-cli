package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const forecastTemplate = `Provide a detailed %d-day weather forecast for %s with specific activity suggestions.

=== RESPONSE REQUIREMENTS ===
* Plain text ONLY - NO markdown, code blocks, or JSON
* Optimized for terminal display (max 140 chars/line)

=== WEATHER FORMAT ===
Location: %[2]s

[Day 1]
📅 Tuesday, June 4 ☀️
🌡 18°C to 24°C | ☔ 10%%%% rain | 💨 10km/h wind
☀️ UV Index: 6 (High)
👕 Wear: Light clothing, sunglasses
🎡 Activity 1: Tempelhofer Feld sunset picnic
🏛 Activity 2: Pergamon Museum new exhibition
🍻 Activity 3: Rooftop bar at Klunkerkranich

=== RULES ===
1. For EACH of %[1]d days:
   - Date with emoji
   - Temperature range with unit
   - Precipitation chance (use %%%%)
   - Wind speed
   - UV index if available
   - MAX 3 specific activities
   - Practical clothing/items advice

2. Activities MUST:
   - Name actual places in %[2]s
   - Match weather conditions
   - Include both indoor/outdoor options
   - Mention current events if available

3. Language: %[3]s`

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type         string        `json:"type"`
	Text         string        `json:"text,omitempty"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatReqBody struct {
	Model    string    `json:"model"`
	Stream   bool      `json:"stream"`
	Messages []Message `json:"messages"`
}

func main() {
	url := "https://openrouter.ai/api/v1/chat/completions"
	apiKey := "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"

	// Get command line arguments
	args := os.Args[1:]

	// Set userInput with fallback to default
	var userInput string
	if len(args) > 0 {
		userInput = args[0]
	} else {
		userInput = "Berlin"
	}

	// Set language with fallback to default
	var language string
	if len(args) > 1 {
		language = args[1]
	} else {
		language = "English"
	}

	// Set maxDays with fallback to default
	var days int
	if len(args) > 2 {
		days, _ = strconv.Atoi(args[2])
	} else {
		days = 3
	}

	// fmt.Println("Enter the city or country you're interested in (default Berlin :)):")
	// var userInput string = "Berlin"
	// fmt.Scanln(&userInput)

	// fmt.Println("Response language (default English):")
	// var language string = "English"
	// fmt.Scanln(&language)

	// Conversation between assistant and user
	messages := []Message{
		{
			Role: "user",
			Content: []Content{
				{
					Type: "text",
					Text: fmt.Sprintf(
						forecastTemplate,
						days,      // %[1]d → days (integer)
						userInput, // %[2]s → location (string)
						language,  // %[3]s → language (string)
					),
				},
			},
		},
	}

	chatRequestBody := ChatReqBody{
		Model:    "google/gemini-2.0-flash-thinking-exp-1219:free",
		Messages: messages,
		Stream:   true,
	}

	requestBody, err := json.Marshal(chatRequestBody)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("Thinking...")

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error: Status %d, Response: %s\n", resp.StatusCode, string(body))
		return
	}

	// Read the response line by line
	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			data = strings.TrimSpace(data)
			if data == "[DONE]" {
				break
			}

			var event map[string]interface{}
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				fmt.Println("Error parsing event:", err)
				continue
			}

			if choices, ok := event["choices"].([]interface{}); ok && len(choices) > 0 {
				firstChoice := choices[0].(map[string]interface{})
				if delta, ok := firstChoice["delta"].(map[string]interface{}); ok {
					if content, ok := delta["content"].(string); ok {
						fmt.Print(content) // Print streamed content
						fullResponse.WriteString(content)
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Print the final formatted response
	// fmt.Println("\n\nFinal Weather Forecast:")
	// fmt.Println(fullResponse.String())
}
