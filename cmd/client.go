package cmd

import (
	"github.com/pojntfx/godibs/pkg/starters"
	"github.com/spf13/cobra"
)

var (
	GIT_UP_BASE_URL string

	PIPELINE_UP_DIR_SRC       string
	PIPELINE_UP_DIR_PUSH      string
	PIPELINE_UP_DIR_WATCH     string
	PIPELINE_UP_FILE_MOD      string
	PIPELINE_UP_BUILD_COMMAND string
	PIPELINE_UP_TEST_COMMAND  string
	PIPELINE_UP_START_COMMAND string
	PIPELINE_UP_REGEX_IGNORE  string

	PIPELINE_DOWN_MODULES     string
	PIPELINE_DOWN_DIR_MODULES string
)

const (
	GIT_UP_COMMIT_MESSAGE = "up_synced"
	GIT_UP_REMOTE_NAME    = "godibs-sync"
	GIT_UP_USER_NAME      = "godibs-syncer"
	GIT_UP_USER_EMAIL     = "godibs-syncer@pojtinger.space"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start the client",
	Run: func(cmd *cobra.Command, args []string) {
		client := starters.Client{
			PipelineUpFileMod:         PIPELINE_UP_FILE_MOD,
			PipelineDownModules:       PIPELINE_DOWN_MODULES,
			PipelineDownDirModules:    PIPELINE_DOWN_DIR_MODULES,
			RedisUrl:                  REDIS_URL,
			RedisPrefix:               REDIS_PREFIX,
			RedisSuffixUpRegistered:   REDIS_SUFFIX_UP_REGISTERED,
			RedisSuffixUpUnRegistered: REDIS_SUFFIX_UP_UNREGISTERED,
			GitUpRemoteName:           GIT_UP_REMOTE_NAME,
			GitUpBaseURL:              GIT_UP_BASE_URL,
			GitUpUserName:             GIT_UP_USER_NAME,
			GitUpUserEmail:            GIT_UP_USER_EMAIL,
			GitUpCommitMessage:        GIT_UP_COMMIT_MESSAGE,
			PipelineUpTestCommand:     PIPELINE_UP_TEST_COMMAND,
			RedisSuffixUpTested:       REDIS_SUFFIX_UP_TESTED,
			PipelineUpBuildCommand:    PIPELINE_UP_BUILD_COMMAND,
			RedisSuffixUpBuilt:        REDIS_SUFFIX_UP_BUILT,
			PipelineUpStartCommand:    PIPELINE_UP_START_COMMAND,
			RedisSuffixUpStarted:      REDIS_SUFFIX_UP_STARTED,
			RedisSuffixUpPushed:       REDIS_SUFFIX_UP_PUSHED,
			PipelineUpDirSrc:          PIPELINE_UP_DIR_SRC,
			PipelineUpDirPush:         PIPELINE_UP_DIR_PUSH,
			PipelineUpDirWatch:        PIPELINE_UP_DIR_WATCH,
			PipelineUpRegexIgnore:     PIPELINE_UP_REGEX_IGNORE,
		}

		client.Start()
	},
}

func init() {
	clientCmd.PersistentFlags().StringVar(&GIT_UP_BASE_URL, "git-base-url", "http://localhost:25000/repos", "Base URL of the sync remote")

	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_DIR_SRC, "dir-src", ".", "Directory in which the source code of the module to push resides")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_DIR_PUSH, "dir-push", "/tmp/.push", "Temporary directory to put the module into before pushing")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_DIR_WATCH, "dir-watch", ".", "Directory to watch for changes")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_FILE_MOD, "modules-file", "go.mod", "Go module file of the module to push")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_BUILD_COMMAND, "cmd-build", "go build ./...", "Command to run to build the module")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_TEST_COMMAND, "cmd-test", "go test ./...", "Command to run to test the module")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_START_COMMAND, "cmd-start", "go run main.go", "Command to run to start the module")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_REGEX_IGNORE, "regex-ignore", "*.pb.go", "Regular expression files to ignore")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_DOWN_MODULES, "modules-pull", "", "Comma-seperated list of the names of the modules to pull")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_DOWN_DIR_MODULES, "dir-pull", "/tmp/modules", "Directory to pull the modules to")

	rootCmd.AddCommand(clientCmd)
}
