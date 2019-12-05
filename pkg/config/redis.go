package config

import "os"

var (
	REDIS_URL    = os.Getenv("REDIS_URL")
	REDIS_PREFIX = os.Getenv("REDIS_PREFIX")
)

const (
	REDIS_SUFFIX_UP_BUILT        = "up_built"
	REDIS_SUFFIX_UP_TESTED       = "up_tested"
	REDIS_SUFFIX_UP_STARTED      = "up_started"
	REDIS_SUFFIX_UP_REGISTERED   = "up_registered"
	REDIS_SUFFIX_UP_UNREGISTERED = "up_unregistered"
	REDIS_SUFFIX_UP_PUSHED       = "up_pushed"

	REDIS_SUFFIX_DOWN_DOWNLOADED = "down_downloaded"
)
