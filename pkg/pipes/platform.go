package pipes

type Platform struct {
	Platform string
	Assets   Assets
	Tests    struct {
		Unit struct {
			Lang Build
		}
		Integration struct {
			Lang   Build
			Image  Build
			Assets Build
		}
	}
}
