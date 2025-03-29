# Installation Guide

This document provides comprehensive installation and configuration instructions for integrating the Strava activity heatmap into your GitHub profile.

## Prerequisites

Before beginning the installation process, ensure you have:

- A GitHub account with a profile README repository (`<username>/<username>`)
- An active Strava account containing activity data
- Basic familiarity with GitHub Actions and repository settings

## Installation Process

### Step 1: Repository Setup

1. **Fork the StravaGraph Repository**

   - Navigate to the main repository page
   - Click the "Fork" button in the upper right corner
   - Wait for GitHub to create your personal copy of the repository

2. **Clone Your Fork** (Optional for local development)
   ```bash
   git clone https://github.com/your-username/StravaGraph.git
   cd StravaGraph
   ```

### Step 2: Strava API Configuration

1. **Create a Strava API Application**

   - Visit [Strava API Settings](https://www.strava.com/settings/api)
   - Click "Create & Manage Your App"
   - Complete the application form with the following details:
     - **Application Name**: GitHub Activity Heatmap (or preferred name)
     - **Category**: Other
     - **Website**: Your GitHub profile URL (e.g., `https://github.com/your-username`)
     - **Authorization Callback Domain**: `localhost`

2. **Obtain API Credentials**
   - After creating the application, note your **Client ID** and **Client Secret**
   - These credentials will be required in subsequent steps

### Step 3: Authentication Setup

1. **Install Dependencies** (for local development)

   ```bash
   # If working with the code locally
   go mod download
   ```

2. **Configure Strava API Credentials**

   #### Option A: Using a .env File (Recommended for local development)

   ```bash
   # Create and populate .env file
   cp .env.example .env
   # Edit .env with your credentials
   ```

   #### Option B: Using Environment Variables

   ```bash
   export STRAVA_CLIENT_ID=your_client_id
   export STRAVA_CLIENT_SECRET=your_client_secret
   ```

3. **Generate Refresh Token**

   ```bash
   go run ./cmd/strava-heatmap/main.go -auth
   ```

4. **Follow Authentication Instructions**
   - You'll be directed to a Strava authorization page
   - After authorizing, you'll receive a code in the redirect URL
   - Exchange this code for a refresh token using the provided curl command
   - Save the refresh token value from the response
   - Add the refresh token to your .env file if using that method

### Step 4: GitHub Configuration

1. **Create GitHub Personal Access Token (PAT)**

   - Navigate to GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Click "Generate new token" → "Generate new token (classic)"
   - Configure the token:
     - **Name**: StravaGraph
     - **Expiration**: Select appropriate duration
     - **Scopes**: Select `repo` (Full control of private repositories)
   - Click "Generate token"
   - **IMPORTANT**: Copy the generated token immediately; it will not be shown again

2. **Configure Repository Secrets**
   - Go to your forked repository's settings
   - Navigate to "Secrets and variables" → "Actions"
   - Add the following repository secrets:
     - `STRAVA_CLIENT_ID`: Your Strava API client ID
     - `STRAVA_CLIENT_SECRET`: Your Strava API client secret
     - `STRAVA_REFRESH_TOKEN`: Your Strava refresh token
     - `PAT`: Your GitHub Personal Access Token

### Step 5: GitHub Profile README Integration

1. **Create Profile Repository** (if needed)

   - Create a new repository with the same name as your GitHub username
   - Initialize it with a README.md file

2. **Add Integration Markers**

   - Edit your profile README.md
   - Add the following markers where you want the heatmap to appear:

   ```markdown
   ## My Strava Activity

   <!-- STRAVA-HEATMAP-START -->
   <!-- STRAVA-HEATMAP-END -->
   ```

3. **Configure Workflow Repository Target**
   - In your forked StravaGraph repository:
   - Locate the `.github/workflows/update-heatmap.yml` file
   - Update the `repository` parameter in the "Checkout profile repository" step:
   ```yaml
   repository: your-username/your-username # Replace with your GitHub username
   ```

### Step 6: Customization

Customize the appearance and behavior of your heatmap by editing the `config.json` file:

```json
{
	"activityTypes": ["Run", "Ride", "Swim", "Hike", "WeightTraining"],
	"metricType": "distance",
	"colorScheme": "strava",
	"showStats": false,
	"dateRange": "1year",
	"cellSize": 10,
	"includePRs": true,
	"darkModeSupport": true,
	"weekStart": "Monday",
	"timeZone": "UTC"
}
```

For a complete reference of all configuration options, see [examples/config.customized.json](./examples/config.customized.json).

### Step 7: Workflow Execution

1. **Trigger Initial Update**

   - Go to the "Actions" tab in your forked repository
   - Select the "Update Strava Heatmap" workflow
   - Click "Run workflow" and confirm

2. **Verify Implementation**
   - Once the workflow completes successfully, visit your GitHub profile
   - Confirm the heatmap appears correctly between the markers
   - The workflow will subsequently run automatically according to the configured schedule (daily at midnight UTC)

## Maintenance

### Token Refresh Handling

- The Strava refresh token is designed for long-term use and should remain valid indefinitely
- The system automatically handles access token refreshes when needed
- You will only need to generate a new refresh token if:
  - You explicitly revoke access for the application
  - You change your Strava account password
  - Strava modifies their security policies

### GitHub PAT Renewal

- Monitor your PAT expiration date if you set a limited duration
- Generate a new token before expiration and update the `PAT` secret in your repository

## Troubleshooting

### Authentication Issues

If encountering authentication errors:

1. **Verify Strava API Credentials**

   - Confirm Client ID and Client Secret are correct
   - Regenerate refresh token if necessary

2. **Test Configuration**
   ```bash
   go run ./cmd/strava-heatmap/main.go -test
   ```

### GitHub Actions Workflow Failures

For issues with the GitHub Actions workflow:

1. **Repository Access Errors (403)**

   - Verify your PAT has the correct `repo` scope
   - Ensure the PAT has not expired
   - Confirm the PAT is correctly stored as a repository secret

2. **Strava API Errors**

   - Check Actions log for specific error messages
   - Verify all required Strava API secrets are correctly configured
   - Ensure the refresh token has the necessary `activity:read_all` scope

3. **README Marker Issues**
   - Confirm your profile README contains both marker comments exactly as shown
   - Ensure there are no typos in the marker syntax

### Advanced Debugging

For more insight into issues:

1. Enable debug mode in `config.json`:

   ```json
   "debug": true
   ```

2. Examine the workflow run logs in the Actions tab of your repository

If you need further assistance, please open an issue in the repository with a description of your problem and relevant logs (with sensitive information redacted).

---

_This installation guide was last updated: March 2025_
