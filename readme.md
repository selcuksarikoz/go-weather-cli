# Weather Forecast CLI

A command-line tool that provides detailed multi-day weather forecasts with location-specific activity suggestions, powered by OpenRouter.ai's API.

## Features

- 🌦️ Get detailed weather forecasts for any location
- 📅 Customizable forecast duration (1-7 days)
- 🌍 Supports multiple languages
- 🎡 Includes practical activity suggestions based on weather conditions
- � Terminal-optimized plain text output
- ⚡ Streamed response for better user experience

## Installation (Global)

```bash
$ sudo cp weather /usr/local/bin/  # Now use 'weather' anywhere
```

```bash
$ weather [location] [language] [days]
```

## Usage

```bash
./weather [location] [language] [days]

Location: Paris

[Day 1]
📅 Tuesday, June 4 ☀️
🌡 18°C to 24°C | ☔ 10% rain | 💨 10km/h wind
☀️ UV Index: 6 (High)
👕 Wear: Light clothing, sunglasses
🎡 Activity 1: Seine River cruise
🏛 Activity 2: Louvre Museum visit
🍻 Activity 3: Montmartre café tour
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

MIT License - see [LICENSE](LICENSE)
