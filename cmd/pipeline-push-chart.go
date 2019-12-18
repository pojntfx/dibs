package cmd

import (
	"github.com/google/uuid"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var PipelinePushChartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Push chart",
	Run: func(cmd *cobra.Command, args []string) {
		id := uuid.New().String()
		pushDir := append([]string{os.TempDir()}, "dibs", "pushHelm", id)

		if output, err := Dibs.PushHelmChart(viper.GetString(PlatformKey), viper.GetString(PushChartGitUserNameKey), viper.GetString(PushChartGitUserEmailKey), viper.GetString(PushChartGitCommitMessageKey), viper.GetString(PushChartGithubUserNameKey), viper.GetString(PushChartGithubTokenKey), viper.GetString(PushChartGithubRepoNameKey), viper.GetString(PushChartGitRepoURLKey), viper.GetString(PushChartGithubPagesURLKey), pushDir); err != nil {
			utils.PipeLogErrorFatalNonPlatformSpecific("Could not Push chart", err, output)
		}
	},
}

func init() {
	var (
		gitUserName      string
		gitUserEmail     string
		gitCommitMessage string
		gitRepoURL       string

		githubUserName string
		githubToken    string
		githubRepoName string
		githubPagesURL string

		gitUserNameFlag      = strings.Replace(strings.Replace(PushChartGitUserNameKey, PushChartKeyPrefix, "", -1), "_", "-", -1)
		gitUserEmailFlag     = strings.Replace(strings.Replace(PushChartGitUserEmailKey, PushChartKeyPrefix, "", -1), "_", "-", -1)
		gitCommitMessageFlag = strings.Replace(strings.Replace(PushChartGitCommitMessageKey, PushChartKeyPrefix, "", -1), "_", "-", -1)
		gitRepoURLFlag       = strings.Replace(strings.Replace(PushChartGitRepoURLKey, PushChartKeyPrefix, "", -1), "_", "-", -1)

		githubUserNameFlag = strings.Replace(strings.Replace(PushChartGithubUserNameKey, PushChartKeyPrefix, "", -1), "_", "-", -1)
		githubTokenFlag    = strings.Replace(strings.Replace(PushChartGithubTokenKey, PushChartKeyPrefix, "", -1), "_", "-", -1)
		githubRepoNameFlag = strings.Replace(strings.Replace(PushChartGithubRepoNameKey, PushChartKeyPrefix, "", -1), "_", "-", -1)
		githubPagesURLFlag = strings.Replace(strings.Replace(PushChartGithubPagesURLKey, PushChartKeyPrefix, "", -1), "_", "-", -1)
	)

	PipelinePushChartCmd.PersistentFlags().StringVar(&gitUserName, gitUserNameFlag, "user", "Git username for the charts repo")
	PipelinePushChartCmd.PersistentFlags().StringVar(&gitUserEmail, gitUserEmailFlag, "user@example.com", "Git user email for the charts repo")
	PipelinePushChartCmd.PersistentFlags().StringVar(&gitCommitMessage, gitCommitMessageFlag, "chore: Update Helm charts", "Git commit message for the charts repo")
	PipelinePushChartCmd.PersistentFlags().StringVar(&gitRepoURL, gitRepoURLFlag, "https://github.com/pojntfx/charts.git", "URL of the charts repo (must be HTTP or HTTPS, not SSH)")

	PipelinePushChartCmd.PersistentFlags().StringVar(&githubUserName, githubUserNameFlag, "user", "Github username")
	PipelinePushChartCmd.PersistentFlags().StringVarP(&githubToken, githubTokenFlag, "t", "1234", "GitHub personal access token")
	PipelinePushChartCmd.PersistentFlags().StringVar(&githubRepoName, githubRepoNameFlag, "charts", "Slug of the GitHub repo to use as the charts repo (don't include the username!)")
	PipelinePushChartCmd.PersistentFlags().StringVar(&githubPagesURL, githubPagesURLFlag, "https://pojntfx.github.io/charts/", "URL of the GitHub pages site to use for the charts repo")

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(PushChartGitUserNameKey, PipelinePushChartCmd.PersistentFlags().Lookup(gitUserNameFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(PushChartGitUserEmailKey, PipelinePushChartCmd.PersistentFlags().Lookup(gitUserEmailFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(PushChartGitCommitMessageKey, PipelinePushChartCmd.PersistentFlags().Lookup(gitCommitMessageFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(PushChartGitRepoURLKey, PipelinePushChartCmd.PersistentFlags().Lookup(gitRepoURLFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}

	if err := viper.BindPFlag(PushChartGithubUserNameKey, PipelinePushChartCmd.PersistentFlags().Lookup(githubUserNameFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(PushChartGithubTokenKey, PipelinePushChartCmd.PersistentFlags().Lookup(githubTokenFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(PushChartGithubRepoNameKey, PipelinePushChartCmd.PersistentFlags().Lookup(githubRepoNameFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(PushChartGithubPagesURLKey, PipelinePushChartCmd.PersistentFlags().Lookup(githubPagesURLFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}

	viper.AutomaticEnv()

	PipelinePushCmd.AddCommand(PipelinePushChartCmd)
}
