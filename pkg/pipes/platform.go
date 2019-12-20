package pipes

// Platform is a target platform like linux/amd64
type Platform struct {
	Platform      string // The Docker platform string (i.e. linux/amd64)
	ChartProfiles ChartProfiles
	Assets        Assets
	Starters      struct {
		Lang Build
	}
	Tests struct { // Tests are the tests for the specific platform
		Unit struct { // Unit tests for the platform
			Lang Build // Unit test using the toolchain
		}
		Integration struct { // Integration tests for the platform
			Lang   Build // Integration test using the toolchain
			Image  Build // Integration test using the Docker image
			Assets Build // Integration test using the assets
			Chart  Build // Integration test using the Helm chart
		}
	}
}
