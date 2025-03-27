# Strava-GitHub Heatmap

A GitHub profile addon that displays your Strava activity as a contribution-style heatmap directly in your GitHub profile README.

## Overview

This tool fetches your Strava activity data and generates a GitHub-style contribution heatmap that can be embedded in your GitHub profile README. It shows your workout frequency and intensity over time, similar to how GitHub visualizes code contributions. The heatmap's color intensity varies based on your activity metrics, giving visitors to your profile a visual representation of your athletic consistency and commitment.

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

## Setup

### Prerequisites

- A GitHub account with a profile README repository (`<username>/<username>`)
- A Strava account with activities
- A Strava API application (for authentication)

### Installation

1. Fork this repository
2. Create a Strava API application at https://www.strava.com/settings/api
   - Set the "Authorization Callback Domain" to `localhost`
3. Add the following secrets to your repository settings (Settings > Secrets and variables > Actions):
   - `STRAVA_CLIENT_ID`: Your Strava API client ID
   - `STRAVA_CLIENT_SECRET`: Your Strava API client secret
   - `STRAVA_REFRESH_TOKEN`: Your Strava refresh token (guide below)
   - `GITHUB_TOKEN`: Automatically provided by GitHub Actions
4. Add the following to your GitHub profile README.md:

```markdown
## My Strava Activity
<!-- STRAVA-HEATMAP-START -->
<!-- STRAVA-HEATMAP-END -->
```

5. Enable GitHub Actions on your fork
6. Run the GitHub Action manually for the first time (Actions tab > Update Strava Heatmap > Run workflow)

### Getting Your Strava Refresh Token

Run the tool with the auth command:

```bash
export STRAVA_CLIENT_ID=your_client_id
export STRAVA_CLIENT_SECRET=your_client_secret
go run ./cmd/strava-heatmap/main.go -auth
```

Follow the instructions to obtain your refresh token.

## Customization

Edit the `config.json` file to customize your heatmap. See the [CLAUDE.md](CLAUDE.md) file for detailed customization options.

## Development

### Building

```bash
go build -o strava-heatmap ./cmd/strava-heatmap
```

### Testing Configuration

```bash
./strava-heatmap -test
```

### Generating SVG Without Updating README

```bash
./strava-heatmap -generate > heatmap.svg
```

### Updating README Manually

```bash
./strava-heatmap -update
```

## License

MIT

---

*Powered by your sweat, built with Go, inspired by GitHub's contribution graph.*

## My Strava Activity
<!-- STRAVA-HEATMAP-START -->
<!-- STRAVA-HEATMAP-END -->