package pipes

type Platform struct {
	Platform string
	Binary   Binary
	Tests    struct {
		Unit struct {
			Lang Build
		}
		Integration struct {
			Lang   Build
			Image  Build
			Binary Build
		}
	}
}
