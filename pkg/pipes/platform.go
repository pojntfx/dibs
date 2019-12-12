package pipes

type Platform struct {
	Platform string
	Binary   Binary
	Tests    struct {
		Unit        Build
		Integration struct {
			Lang   Build
			Image  Build
			Binary Build
		}
	}
}
