# Installation Guide

This guide will walk you through the process of setting up the Strava-GitHub Heatmap for your GitHub profile.

## Prerequisites

- A GitHub account with a profile README repository (`<username>/<username>`)
- A Strava account with activities
- A Strava API application (for authentication)

## Step 1: Fork This Repository

1. Click the "Fork" button at the top right of this repository.
2. This will create a copy of the repository under your GitHub account.

## Step 2: Create a Strava API Application

1. Go to [Strava API Settings](https://www.strava.com/settings/api)
2. Create a new application with the following details:
   - **Application Name**: GitHub Activity Heatmap (or any name you prefer)
   - **Category**: Other
   - **Website**: Your GitHub profile URL (e.g., `https://github.com/yourusername`)
   - **Authorization Callback Domain**: `localhost`
3. After creating the application, note down your **Client ID** and **Client Secret**.

## Step 3: Get Your Strava Refresh Token

1. Install the tool or run it directly using Go:
   ```
   go run ./cmd/strava-heatmap/main.go -auth
   ```
   
2. Set the environment variables with your Strava API credentials:
   ```
   export STRAVA_CLIENT_ID=your_client_id
   export STRAVA_CLIENT_SECRET=your_client_secret
   ```

3. Follow the instructions displayed to authorize the application and obtain your refresh token.
   - You'll be directed to a Strava authorization page
   - After authorization, you'll get a code that can be exchanged for a refresh token
   - Use the provided curl command to exchange the code
   - Copy the refresh token from the response

## Step 4: Configure GitHub Repository Secrets

1. Go to your forked repository's settings.
2. Navigate to "Secrets and variables" > "Actions".
3. Add the following repository secrets:
   - `STRAVA_CLIENT_ID`: Your Strava API client ID
   - `STRAVA_CLIENT_SECRET`: Your Strava API client secret
   - `STRAVA_REFRESH_TOKEN`: The refresh token you obtained in Step 3

## Step 5: Set Up Your GitHub Profile README

1. Create a GitHub profile README repository if you don't have one already:
   - Create a new repository with the same name as your GitHub username
   - Add a README.md file to this repository

2. Add the following markers to your README where you want the Strava heatmap to appear:
   ```markdown
   ## My Strava Activity
   <!-- STRAVA-HEATMAP-START -->
   <!-- STRAVA-HEATMAP-END -->
   ```

## Step 6: Customize Your Heatmap (Optional)

Edit the `config.json` file in your forked repository to customize the appearance and behavior of your heatmap:

```json
{
  "activityTypes": ["Run", "Ride", "Swim", "Hike", "WeightTraining"],
  "metricType": "distance",
  "colorScheme": "github",
  "showStats": true,
  "dateRange": "1year",
  "cellSize": 11,
  "includePRs": true,
  "includeLocationHeatmap": false,
  "darkModeSupport": true,
  "weekStart": "Monday",
  "language": "en",
  "timeZone": "UTC"
}
```

See the [README.md](README.md) for detailed information about the configuration options.

## Step 7: Run the GitHub Action

1. Go to the "Actions" tab in your forked repository.
2. Click on the "Update Strava Heatmap" workflow.
3. Click "Run workflow" to manually trigger the action.

The action will:
- Build the tool
- Fetch your Strava activities
- Generate the heatmap SVG
- Update your profile README
- Commit and push the changes

The workflow is also scheduled to run daily at midnight UTC to keep your heatmap updated.

## Troubleshooting

- **Authentication Issues**: If you encounter authentication errors, try regenerating your refresh token.
- **Missing Activities**: Ensure you've granted the appropriate permissions when authorizing the Strava application.
- **Workflow Failures**: Check the GitHub Actions logs for detailed error messages.

If you need further assistance, please open an issue in the repository.