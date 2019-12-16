package pipes

type Platform struct {
	Platform      string
	ChartProfiles ChartProfiles
	Assets        Assets
	Tests         struct {
		Unit struct {
			Lang Build
		}
		Integration struct {
			Lang   Build
			Image  Build
			Assets Build
			Chart  Build
		}
	}
}
