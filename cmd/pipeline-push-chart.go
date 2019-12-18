package cmd

import (
	"github.com/google/uuid"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"os"
	"strings"
)

var PipelinePushChartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Push chart",
	Run: func(cmd *cobra.Command, args []string) {
		id := uuid.New().String()
		pushDir := append([]string{os.TempDir()}, "dibs", "pushHelm", id)

		if output, err := Dibs.PushHelmChart(viper.GetString(PlatformKey), viper.GetString(GitUserNameKey), viper.GetString(GitUserEmailKey), viper.GetString(GitCommitMessageKey), viper.GetString(GithubUserNameKey), viper.GetString(GithubTokenKey), viper.GetString(GithubRepoNameKey), viper.GetString(GitRepoURLKey), viper.GetString(GithubPagesURLKey), pushDir); err != nil {
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

		gitUserNameFlag      = strings.Replace(GitUserNameKey, "_", "-", -1)
		gitUserEmailFlag     = strings.Replace(GitUserEmailKey, "_", "-", -1)
		gitCommitMessageFlag = strings.Replace(GitCommitMessageKey, "_", "-", -1)
		gitRepoURLFlag       = strings.Replace(GitRepoURLKey, "_", "-", -1)

		githubUserNameFlag = strings.Replace(GithubUserNameKey, "_", "-", -1)
		githubTokenFlag    = strings.Replace(GithubTokenKey, "_", "-", -1)
		githubRepoNameFlag = strings.Replace(GithubRepoNameKey, "_", "-", -1)
		githubPagesURLFlag = strings.Replace(GithubPagesURLKey, "_", "-", -1)
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

	if err := viper.BindPFlag(GitUserNameKey, PipelinePushChartCmd.PersistentFlags().Lookup(gitUserNameFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(GitUserEmailKey, PipelinePushChartCmd.PersistentFlags().Lookup(gitUserEmailFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(GitCommitMessageKey, PipelinePushChartCmd.PersistentFlags().Lookup(gitCommitMessageFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(GitRepoURLKey, PipelinePushChartCmd.PersistentFlags().Lookup(gitRepoURLFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(GithubUserNameKey, PipelinePushChartCmd.PersistentFlags().Lookup(githubUserNameFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(GithubTokenKey, PipelinePushChartCmd.PersistentFlags().Lookup(githubTokenFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(GithubRepoNameKey, PipelinePushChartCmd.PersistentFlags().Lookup(githubRepoNameFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(GithubPagesURLKey, PipelinePushChartCmd.PersistentFlags().Lookup(githubPagesURLFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	viper.AutomaticEnv()

	PipelinePushCmd.AddCommand(PipelinePushChartCmd)
}
