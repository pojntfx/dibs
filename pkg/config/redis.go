package config

import "os"

var (
	REDIS_URL            = os.Getenv("REDIS_URL")
	REDIS_CHANNEL_PREFIX = os.Getenv("REDIS_CHANNEL_PREFIX")
)

const (
	REDIS_CHANNEL_MODULE_BUILT        = "module_built"
	REDIS_CHANNEL_MODULE_TESTED       = "module_tested"
	REDIS_CHANNEL_MODULE_STARTED      = "module_started"
	REDIS_CHANNEL_MODULE_REGISTERED   = "module_registered"
	REDIS_CHANNEL_MODULE_UNREGISTERED = "module_unregistered"
	REDIS_CHANNEL_MODULE_PUSHED       = "module_pushed"
)
