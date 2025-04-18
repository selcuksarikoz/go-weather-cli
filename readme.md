````markdown
# Weather Forecast CLI

A command-line tool that provides detailed multi-day weather forecasts with location-specific activity suggestions, powered by OpenRouter.ai's API.

## Features

- 🌦️ Get detailed weather forecasts for any location
- 📅 Customizable forecast duration (1-7 days)
- 🌍 Supports multiple languages
- 🎡 Includes practical activity suggestions based on weather conditions
- � Terminal-optimized plain text output
- ⚡ Streamed response for better user experience

## Installation

1. Ensure you have Go installed (version 1.20+ recommended)
2. Clone this repository or download the source code
3. Build the application:
   ```bash
   go build -o weather
   ```
````

## Usage

```bash
./weather [location] [language] [days]
```

### Arguments

- `location`: City or country name (default: "Berlin")
- `language`: Response language (default: "English")
- `days`: Number of forecast days (1-7, default: 3)

### Examples

1. Default forecast for Berlin:

   ```bash
   ./weather
   ```

2. 5-day forecast for Paris in French:

   ```bash
   ./weather Paris French 5
   ```

3. Weekend forecast for Tokyo in Japanese:
   ```bash
   ./weather Tokyo Japanese 2
   ```

## Configuration

You'll need an OpenRouter.ai API key:

1. Sign up at [OpenRouter.ai](https://openrouter.ai/)
2. Get your API key from the dashboard
3. Replace the `apiKey` constant in `main.go` with your key

## Response Format

The forecast includes for each day:

- 📅 Date and weather emoji
- 🌡 Temperature range
- ☔ Precipitation chance
- 💨 Wind speed
- ☀️ UV Index (when available)
- 👕 Clothing recommendations
- 🎡 3 specific activity suggestions

## Requirements

- Go 1.20+
- Internet connection
- OpenRouter.ai API key

## Limitations

- Free tier API may have rate limits
- Accuracy depends on the underlying model's knowledge
- Activity suggestions may not always reflect current events

## License

MIT License - see [LICENSE](LICENSE) file (create one if needed)

```

You might want to add these additional files:

1. `LICENSE` - Add an appropriate open-source license
2. `.gitignore` - For Go projects
3. `go.mod` - If you want to make it a proper Go module

The README provides clear installation instructions, usage examples, and explains the features in a user-friendly way. You can customize it further based on your specific needs or additional features you plan to add.
```
