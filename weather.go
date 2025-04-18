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
	"time"
)

const forecastTemplate = `Generate a precise %d-day weather forecast for %s starting from %s with hyper-localized activity recommendations.

=== RESPONSE REQUIREMENTS ===
* STRICT plain text format (NO markdown, JSON, or code blocks)
* Terminal-optimized (max 138 characters/line)
* Use ONLY these weather emojis: ☀️🌤️⛅🌥️☁️🌦️🌧️⛈️🌨️💨🌫️
* Include real-time weather anomalies if present
* Date calculations MUST start from %[4]s (YYYY-MM-DD)

=== WEATHER DATA FORMAT ===
Location: %[2]s | 📍 [Google Maps URL: https://maps.google.com/?q=%[2]s]
🌿 Pollen Alert: [Current pollen types & intensity]

------------------------------------------------------------------
[Day X]
📅 [Weekday, Month DD] (Calculated from %[4]s + X days) [Weather Emoji] 
🌡 [Min]°C to [Max]°C | ☔ [Percip]%% | 💨 [Wind]km/h [Direction]
🌞 UV Index: [Value] ([Level]) | 🌙 Overnight: [Temp]°C
👕 Wear: [Clothing items]
🎩 Fancy Tip: [Humorous fashion suggestion] 
🤣 Local Quirk: [Amusing local dressing observation]
🤧 Allergy: [Pollen advice]

🏙️ Activity 1: [Specific venue/event] ([Indoor/Outdoor])  
📍 [Google Maps: https://maps.google.com/?q=[Venue+Name]]  
🎨 Activity 2: [Current exhibition/performance]  
📍 [Google Maps: https://maps.google.com/?q=[Exhibition+Venue]]  
🍽️ Activity 3: [Weather-appropriate dining]  
📍 [Google Maps: https://maps.google.com/?q=[Restaurant+Name]]  
------------------------------------------------------------------

=== STRICT RULES ===
1. Date calculations MUST begin from %[4]s for all day forecasts
2. Humorous additions MUST:
   - Be culturally appropriate for %[2]s
   - Reference local fashion stereotypes or weather quirks
   - Never exceed 2 lines (max 138 chars each)
3. Activities MUST:  
   - Include VERIFIABLE venues/events in %[2]s  
   - Provide DIRECT Google Maps links  
   - Specify travel method (e.g., "10-min walk from [Station]")  
4. Weather-dependent adjustments REQUIRED
5. Language: %[3]s (with local idioms where funny)
`

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
	apiKey := ""
	currentDate := time.Now()
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
						days,
						userInput,
						language,
						currentDate,
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
