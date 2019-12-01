package config

import "os"

var (
	GIT_BASE_URL = os.Getenv("GIT_BASE_URL")

	GIT_SERVER_REPOS_DIR = os.Getenv("GIT_SERVER_REPOS_DIR")
	GIT_SERVER_HTTP_PORT = os.Getenv("GIT_SERVER_HTTP_PORT")
	GIT_SERVER_HTTP_PATH = os.Getenv("GIT_SERVER_HTTP_PATH")

	GIT_UP_REMOTE_NAME = os.Getenv("GIT_UP_REMOTE_NAME")
	GIT_UP_USER_NAME   = os.Getenv("GIT_UP_USER_NAME")
	GIT_UP_USER_EMAIL  = os.Getenv("GIT_UP_USER_EMAIL")
)

const (
	GIT_UP_COMMIT_MESSAGE = "up_synced"
)
