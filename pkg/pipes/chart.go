package pipes

// Chart is a Helm chart
type Chart struct {
	SrcDir     string   // Source directory of the Helm chart
	DistDir    string   // Directory into which the Helm chart should be placed
	CleanGlobs []string // Array of globs to remove when cleaning
}
