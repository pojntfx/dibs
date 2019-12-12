package utils

type Platform struct {
	Platform string
	Image    Build
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
