# Strava-GitHub Heatmap 🏃‍♂️ 🚴‍♀️ 📊

[![GitHub Stars](https://img.shields.io/github/stars/samuellee/StravaGraph?style=social)](https://github.com/samuellee/StravaGraph/stargazers)
[![GitHub License](https://img.shields.io/github/license/samuellee/StravaGraph)](https://github.com/samuellee/StravaGraph/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/samuellee/StravaGraph)](https://goreportcard.com/report/github.com/samuellee/StravaGraph)
[![Last Release](https://img.shields.io/github/v/release/samuellee/StravaGraph)](https://github.com/samuellee/StravaGraph/releases)

A GitHub profile addon that displays your Strava activity as a contribution-style heatmap directly in your GitHub profile README.

## Overview

This tool fetches your Strava activity data and generates a GitHub-style contribution heatmap that can be embedded in your GitHub profile README. It shows your workout frequency and intensity over time, similar to how GitHub visualizes code contributions. The heatmap's color intensity varies based on your activity metrics, giving visitors to your profile a visual representation of your athletic consistency and commitment.

![Example Strava Heatmap](./examples/example-heatmap.png)

## Features

- 🔥 Displays your Strava activities in a GitHub-style heatmap
- 🏃‍♂️ Supports running, cycling, swimming and other activity types
- 🌈 Color intensity based on activity metrics (configurable):
  - Duration: Time spent exercising
  - Distance: Kilometers/miles covered
  - Elevation: Meters/feet climbed
  - Effort: Calculated metric combining factors
  - Heart rate: Average or max heart rate zones
- 🔄 Auto-updates daily via GitHub Actions
- 🔐 Secure authentication with Strava API (no exposing tokens)
- 📊 Optional statistics summary (weekly, monthly, yearly)
- 🎯 Customizable appearance to match your GitHub profile theme
- 🌍 Activity location heatmap option (with privacy controls)
- 🏆 Achievement highlighting for PRs and significant milestones

## Quick Start

1. **Fork this repository**
2. **Create a Strava API application** at https://www.strava.com/settings/api
3. **Get your Strava refresh token** by following the instructions after running:
   ```bash
   export STRAVA_CLIENT_ID=your_client_id
   export STRAVA_CLIENT_SECRET=your_client_secret
   go run ./cmd/strava-heatmap/main.go -auth
   ```
4. **Add secrets to your repository** (Settings > Secrets and variables > Actions):
   - `STRAVA_CLIENT_ID`: Your Strava API client ID
   - `STRAVA_CLIENT_SECRET`: Your Strava API client secret
   - `STRAVA_REFRESH_TOKEN`: Your Strava refresh token
5. **Add the markers to your GitHub profile README**:
   ```markdown
   ## My Strava Activity
   <!-- STRAVA-HEATMAP-START -->
   <!-- STRAVA-HEATMAP-END -->
   ```
6. **Run the GitHub Action** to update your heatmap

For detailed setup instructions, see the [Installation Guide](./INSTALL.md).

## Usage

### Building the Tool

```bash
go build -o strava-heatmap ./cmd/strava-heatmap
```

### Available Commands

- **Authentication Instructions**: `-auth`
  ```bash
  ./strava-heatmap -auth
  ```
- **Update README with Heatmap**: `-update`
  ```bash
  ./strava-heatmap -update
  ```
- **Generate SVG Without Updating README**: `-generate`
  ```bash
  ./strava-heatmap -generate > heatmap.svg
  ```
- **Test Configuration and Authentication**: `-test`
  ```bash
  ./strava-heatmap -test
  ```

### Customization

Edit the `config.json` file to customize your heatmap:

```json
{
  "activityTypes": ["Run", "Ride", "Swim", "Hike", "WeightTraining"],
  "metricType": "distance",
  "colorScheme": "github",
  "showStats": true,
  "dateRange": "1year",
  "cellSize": 11,
  "includePRs": true,
  "darkModeSupport": true,
  "weekStart": "Monday",
  "timeZone": "UTC"
}
```

See [examples/config.customized.json](./examples/config.customized.json) for a full example with all options.

## Documentation

- [Installation Guide](./INSTALL.md) - Detailed setup instructions
- [API Documentation](./API.md) - Internal API documentation
- [Example Profile](./examples/profile/README.md) - Example GitHub profile README

## Development

### Project Structure

```
/
├── cmd/
│   └── strava-heatmap/
│       └── main.go                 # Main entry point
├── internal/
│   ├── auth/                       # Strava authentication
│   ├── strava/                     # Strava API client
│   ├── processor/                  # Data processing
│   ├── svg/                        # SVG generation
│   ├── github/                     # GitHub integration
│   └── config/                     # Configuration
├── .github/
│   └── workflows/
│       └── update-heatmap.yml      # GitHub Action workflow
├── examples/
│   ├── profile/                    # Example profile README
│   └── config.customized.json      # Example custom config
├── config.json                     # Configuration file
├── README.md                       # This file
├── INSTALL.md                      # Installation guide
└── API.md                          # API documentation
```

### Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code follows the project's style and includes appropriate tests.

## License

MIT

## Acknowledgements

- Thanks to GitHub for the inspiration from their contribution graph
- Thanks to Strava for providing the API
- Icons provided by [Feather Icons](https://feathericons.com/)
- Color schemes inspired by [GitHub](https://github.com) and [Strava](https://www.strava.com)

---

*Powered by your sweat, built with Go, inspired by GitHub's contribution graph.*

## My Strava Activity
<!-- STRAVA-HEATMAP-START -->
<!-- STRAVA-HEATMAP-END -->