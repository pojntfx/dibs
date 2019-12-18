package pipes

// ChartProfiles the Skaffold profiles for the Helm chart
type ChartProfiles struct {
	Production  string // Skaffold profile for production
	Development string // Skaffold profile for development
}
