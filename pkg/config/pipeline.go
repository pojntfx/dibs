package config

import "os"

var (
	SRC_DIR                     = os.Getenv("SRC_DIR")
	PUSH_DIR                    = os.Getenv("PUSH_DIR")
	SYNC_MODULE_PUSH_MOD_FILE   = os.Getenv("SYNC_MODULE_PUSH_MOD_FILE")
	SYNC_MODULE_PUSH_WATCH_GLOB = os.Getenv("SYNC_MODULE_PUSH_WATCH_GLOB")
	COMMAND_BUILD               = os.Getenv("COMMAND_BUILD")
	COMMAND_TEST                = os.Getenv("COMMAND_TEST")
	COMMAND_START               = os.Getenv("COMMAND_START")
)
