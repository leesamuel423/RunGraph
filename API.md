# API Documentation

This document provides detailed information about the internal API of the Strava-GitHub Heatmap tool.

## Modules

### Config Module (`internal/config`)

The config module handles application configuration loading, parsing, and validation.

#### Main Types:

- **Config**: Represents the application configuration.
  ```go
  type Config struct {
      ActivityTypes        []string
      MetricType           string
      ColorScheme          string
      CustomColors         []string
      ShowStats            bool
      StatTypes            []string
      DateRange            string
      CustomDateRange      struct {
          Start string
          End   string
      }
      CellSize              int
      IncludePRs            bool
      IncludeLocationHeatmap bool
      LocationPrivacyRadius int
      DarkModeSupport       bool
      DarkModeColors        []string
      WeekStart             string
      Language              string
      TimeZone              string
      Debug                 bool
  }
  ```

#### Main Functions:

- **LoadConfig(filePath string) (*Config, error)**: Loads configuration from a file.
- **ValidateConfig(config *Config) error**: Validates the configuration values.
- **SaveConfig(config *Config, filePath string) error**: Saves configuration to a file.
- **GetTimeZoneLocation() (*time.Location, error)**: Returns the time.Location for the configured timezone.
- **GetDateRange() (time.Time, time.Time, error)**: Returns the start and end time for the configured date range.

### Authentication Module (`internal/auth`)

The authentication module handles Strava API authentication.

#### Main Types:

- **TokenManager**: Manages Strava API access tokens.
  ```go
  type TokenManager struct {
      ClientID     string
      ClientSecret string
      RefreshToken string
      AccessToken  string
      ExpiresAt    time.Time
  }
  ```

- **TokenResponse**: Represents the response from Strava token endpoint.
  ```go
  type TokenResponse struct {
      TokenType    string
      ExpiresAt    int64
      ExpiresIn    int
      RefreshToken string
      AccessToken  string
      Athlete      struct {
          ID int64
      }
  }
  ```

- **OAuthConfig**: Holds OAuth configuration for Strava.
  ```go
  type OAuthConfig struct {
      ClientID     string
      ClientSecret string
      RedirectURI  string
      Scopes       []string
  }
  ```

#### Main Functions:

- **NewTokenManager(clientID, clientSecret, refreshToken string) *TokenManager**: Creates a new token manager.
- **GetAccessToken() (string, error)**: Returns a valid access token, refreshing if necessary.
- **RefreshAccessToken() error**: Refreshes the Strava access token using the refresh token.
- **GetAuthorizationURL() string**: Returns the URL to redirect the user for authorization.
- **ExchangeCodeForToken(code string) (*TokenResponse, error)**: Exchanges an authorization code for tokens.
- **GetInstructionsForUserAuth(clientID, clientSecret string) string**: Returns instructions for manual token acquisition.

### Strava Module (`internal/strava`)

The Strava module handles interaction with the Strava API.

#### Main Types:

- **TokenManager**: Interface that defines methods for token management.
  ```go
  type TokenManager interface {
      GetAccessToken() (string, error)
      RefreshAccessToken() error
  }
  ```

- **Client**: Handles API communication with Strava.
  ```go
  type Client struct {
      httpClient   *http.Client
      tokenManager TokenManager
      debug        bool
  }
  ```

- **SummaryActivity**: Represents a summary of an activity from Strava API.
  ```go
  type SummaryActivity struct {
      ID              int64
      Name            string
      Distance        float64
      MovingTime      int
      ElapsedTime     int
      TotalElevGain   float64
      Type            string
      StartDate       time.Time
      StartDateLocal  time.Time
      Timezone        string
      // additional fields omitted for brevity
  }
  ```

- **DailyActivity**: Represents aggregated activities for a single day.
  ```go
  type DailyActivity struct {
      Date           time.Time
      Count          int
      TotalDistance  float64
      TotalDuration  int
      TotalElevation float64
      Activities     []int64
      MaxHeartRate   float64
      AvgHeartRate   float64
      HasPR          bool
      Types          map[string]int
  }
  ```

- **HeatmapIntensity**: Represents the intensity level for the heatmap cell.
  ```go
  type HeatmapIntensity int
  
  const (
      None HeatmapIntensity = iota
      Low
      Medium
      High
      VeryHigh
  )
  ```

#### Main Functions:

- **NewClient(tokenManager TokenManager, debug bool) *Client**: Creates a new Strava API client.
- **GetAthlete() (map[string]interface{}, error)**: Gets the authenticated athlete's profile.
- **GetActivities(after, before time.Time, page, perPage int) ([]SummaryActivity, error)**: Retrieves activities for the authenticated athlete.
- **GetAllActivities(after, before time.Time, types []string) ([]SummaryActivity, error)**: Retrieves all activities within the given time range.

### Processor Module (`internal/processor`)

The processor module handles activity data processing and aggregation.

#### Main Types:

- **ActivityAggregator**: Processes and aggregates activity data.
  ```go
  type ActivityAggregator struct {
      Activities []strava.SummaryActivity
      TimeZone   *time.Location
      DailyData  map[string]*strava.DailyActivity
  }
  ```

- **MetricsCalculator**: Calculates activity metrics.
  ```go
  type MetricsCalculator struct {
      DailyData []*strava.DailyActivity
      StartDate time.Time
      EndDate   time.Time
  }
  ```

- **StatsGenerator**: Generates comprehensive statistics.
  ```go
  type StatsGenerator struct {
      DailyData  []*strava.DailyActivity
      StartDate  time.Time
      EndDate    time.Time
      MetricType string
  }
  ```

#### Main Functions:

- **NewActivityAggregator(activities []strava.SummaryActivity, location *time.Location) *ActivityAggregator**: Creates a new activity aggregator.
- **Aggregate() map[string]*strava.DailyActivity**: Processes activities and aggregates them by day.
- **GetOrderedDates(startDate, endDate time.Time) []*strava.DailyActivity**: Returns daily activities ordered by date.
- **CalculateIntensity(metricType string, day *strava.DailyActivity) strava.HeatmapIntensity**: Determines the heat intensity level for a given metric value.
- **CalculateOverallStats() *strava.ActivityStats**: Calculates overall activity statistics.
- **CalculatePeriodStats(periodType string) []*strava.DatePeriodStats**: Calculates statistics for specific time periods.
- **CalculateAverages() map[string]float64**: Calculates average metrics per active day.
- **CalculateEffortScore() float64**: Calculates an overall effort score.
- **GenerateStats() map[string]interface{}**: Generates all statistics for the heatmap.

### SVG Module (`internal/svg`)

The SVG module handles the generation of SVG visualizations.

#### Main Types:

- **Generator**: Handles SVG generation.
  ```go
  type Generator struct {
      Config *config.Config
      Debug  bool
  }
  ```

- **HeatmapData**: Holds all data needed to generate the heatmap.
  ```go
  type HeatmapData struct {
      StartDate       time.Time
      EndDate         time.Time
      Cells           [][]*HeatmapCell
      WeekLabels      []string
      MonthLabels     []struct {
          Month string
          X     int
      }
      ColorTheme      ColorTheme
      DarkModeTheme   ColorTheme
      CellSize        int
      CellSpacing     int
      WeekStart       string
      DarkModeSupport bool
  }
  ```

- **HeatmapCell**: Represents a single cell in the heatmap.
  ```go
  type HeatmapCell struct {
      Date      time.Time
      Intensity strava.HeatmapIntensity
      HasPR     bool
      Count     int
      Tooltip   string
  }
  ```

- **ColorTheme**: Represents a set of colors for the heatmap.
  ```go
  type ColorTheme struct {
      Name   string
      Colors []string
  }
  ```

- **TooltipData**: Holds data needed for a tooltip.
  ```go
  type TooltipData struct {
      Date            time.Time
      ActivityCount   int
      TotalDistance   float64
      TotalDuration   int
      TotalElevation  float64
      ActivityTypes   map[string]int
      HasPR           bool
      CustomFields    map[string]string
  }
  ```

#### Main Functions:

- **NewGenerator(cfg *config.Config) *Generator**: Creates a new SVG generator.
- **GenerateHeatmap(activities []strava.SummaryActivity) (string, error)**: Creates a heatmap SVG from activity data.
- **GenerateLocationHeatmap(activities []strava.SummaryActivity, privacyRadius int) (string, error)**: Creates a heatmap of activity locations.
- **NewHeatmapData(activities []*strava.DailyActivity, startDate, endDate time.Time, ...) *HeatmapData**: Creates a new heatmap data structure.
- **RenderSVG() string**: Generates the SVG for the heatmap with a 7-row layout (one row per day of the week).
- **GetTheme(name string, customColors []string) ColorTheme**: Returns a color theme by name.
- **GetDarkModeTheme(lightTheme ColorTheme, customDarkColors []string) ColorTheme**: Returns the dark mode variant of a color theme.
- **GenerateTooltipSVG(data *TooltipData) string**: Creates an SVG tooltip.

### GitHub Module (`internal/github`)

The GitHub module handles GitHub integration for updating README files and GitHub Actions.

#### Main Types:

- **ReadmeUpdater**: Handles updating the GitHub profile README.
  ```go
  type ReadmeUpdater struct {
      FilePath string
      Debug    bool
  }
  ```

- **ActionsHandler**: Helps with GitHub Actions integration.
  ```go
  type ActionsHandler struct {
      Debug bool
  }
  ```

#### Main Functions:

- **NewReadmeUpdater(filePath string, debug bool) *ReadmeUpdater**: Creates a new README updater.
- **UpdateReadme(svgContent string) error**: Updates the README with the generated SVG.
- **ValidateReadme() (bool, error)**: Checks if the README has the required markers.
- **NewActionsHandler(debug bool) *ActionsHandler**: Creates a new GitHub Actions handler.
- **SetOutput(name, value string) error**: Sets a GitHub Actions output variable.
- **LogError(msg string, err error)**: Logs an error in a GitHub Actions friendly format.
- **LogWarning(msg string)**: Logs a warning in a GitHub Actions friendly format.
- **LogInfo(msg string)**: Logs an info message in a GitHub Actions friendly format.
- **GetEnvWithFallback(key, fallback string) string**: Gets an environment variable with a fallback value.
- **IsRunningInActions() bool**: Checks if the code is running in GitHub Actions.
- **RecordMetric(name string, value interface{})**: Records a metric for the GitHub Action.
- **CreateSummary(content string) error**: Adds content to the GitHub Actions step summary.
- **FormatTimestamp(t time.Time) string**: Formats a timestamp for GitHub Actions logs.

## Command Line Interface

The command line interface is implemented in `cmd/strava-heatmap/main.go` and provides the following commands:

- **-auth**: Generate authentication instructions
- **-update**: Update the heatmap in the README
- **-generate**: Generate SVG without updating README
- **-test**: Test configuration and authentication

## Configuration Schema

The configuration file (`config.json`) follows this schema:

```json
{
  "activityTypes": ["Run", "Ride", "Swim", "Hike", "WeightTraining"],
  "metricType": "distance",
  "colorScheme": "strava",
  "customColors": ["#494950", "#ffd4d1", "#ffad9f", "#fc7566", "#e34a33"],
  "showStats": false,
  "statTypes": ["weekly", "monthly", "yearly"],
  "dateRange": "1year",
  "customDateRange": {
    "start": "2023-01-01",
    "end": "2023-12-31"
  },
  "cellSize": 10,
  "includePRs": true,
  "includeLocationHeatmap": false,
  "locationPrivacyRadius": 500,
  "darkModeSupport": true,
  "darkModeColors": ["#36363c", "#7c2c2a", "#a63b33", "#d64c3b", "#fc7566"],
  "weekStart": "Monday",
  "language": "en",
  "timeZone": "UTC",
  "debug": false
}
```

Valid values:
- **metricType**: "distance", "duration", "elevation", "effort", "heart_rate"
- **colorScheme**: "github", "strava", "blue", "purple", "custom"
- **dateRange**: "1year", "all", "ytd", "custom"
- **weekStart**: "Sunday", "Monday"
- **statTypes**: "weekly", "monthly", "yearly"
