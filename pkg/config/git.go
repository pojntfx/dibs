package config

var (
	GIT_SERVER_REPOS_DIR string
	GIT_SERVER_HTTP_PORT string
	GIT_SERVER_HTTP_PATH string

	GIT_UP_REMOTE_NAME string
	GIT_UP_USER_NAME   string
	GIT_UP_USER_EMAIL  string
	GIT_UP_BASE_URL    string
)

const (
	GIT_UP_COMMIT_MESSAGE = "up_synced"
)
