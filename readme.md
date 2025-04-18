# Weather Forecast CLI

> Instant weather forecasts with activity suggestions (pre-built binary included)

## Quick Start

1. **Download** the `weather` binary
2. **Make executable**:
   ```bash
   chmod +x weather
   ```
3. **Run**:
   ```bash
   ./weather [city] [language] [days]
   ```

## Usage Examples

```bash
./weather              # Berlin forecast (3 days, English)
./weather Tokyo        # Tokyo forecast
./weather Paris French # In French language
./weather "" "" 5     # 5-day forecast (defaults)
```

## Optional: Install Globally

```bash
sudo cp weather /usr/local/bin/  # Now use just 'weather' anywhere
```

> That's it! No Go installation needed - just download and run.
