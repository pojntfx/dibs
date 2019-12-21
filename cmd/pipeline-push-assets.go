package cmd

import (
	"github.com/pojntfx/dibs/pkg/pipes"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var PipelinePushAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Push assets",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		Dibs.RunForPlatformsSerially(platformFromConfig, platformFromConfig == PlatformAll, func(platform pipes.Platform) {
			if output, err := platform.Assets.Push(platform.Platform, strings.Split(viper.GetString(PushAssetsVersionKey), " "), viper.GetString(PushAssetsGitHubTokenKey), viper.GetString(PushAssetsGithubUserNameKey), viper.GetString(PushAssetsGithubRepoNameKey)); err != nil {
				utils.LogErrorFatalPlatformSpecific("Could not push assets", err, platform.Platform, output)
			}
		})
	},
}

func init() {
	var (
		version string

		githubToken    string
		githubRepoName string
		githubUserName string

		versionFlag = strings.Replace(strings.Replace(PushAssetsVersionKey, PushAssetsKeyPrefix, "", -1), "_", "-", -1)

		githubTokenFlag    = strings.Replace(strings.Replace(PushAssetsGitHubTokenKey, PushAssetsKeyPrefix, "", -1), "_", "-", -1)
		githubRepoNameFlag = strings.Replace(strings.Replace(PushAssetsGithubRepoNameKey, PushAssetsKeyPrefix, "", -1), "_", "-", -1)
		githubUserNameFlag = strings.Replace(strings.Replace(PushAssetsGithubUserNameKey, PushAssetsKeyPrefix, "", -1), "_", "-", -1)
	)

	// If we are in a Git repo, get the latest tag and use it as the default version. Ignore errors if we aren't in a Git repo
	pwd, _ := os.Getwd()
	currentGitRepo := &utils.Git{WorkDir: pwd}
	latestGitTag, _ := currentGitRepo.GetLatestTag()

	PipelinePushAssetsCmd.PersistentFlags().StringVarP(&version, versionFlag, "v", latestGitTag, `The version of the asset to deploy (newest Git tag by default) (use "-prerelease <version>" as the version to create a prerelease)`)

	PipelinePushAssetsCmd.PersistentFlags().StringVarP(&githubToken, githubTokenFlag, "t", "1234", "GitHub personal access token")
	PipelinePushAssetsCmd.PersistentFlags().StringVar(&githubRepoName, githubRepoNameFlag, "releases", "Slug of the GitHub repo to push the assets to (don't include the username!)")
	PipelinePushAssetsCmd.PersistentFlags().StringVar(&githubUserName, githubUserNameFlag, "user", "Github username")

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(PushAssetsVersionKey, PipelinePushAssetsCmd.PersistentFlags().Lookup(versionFlag)); err != nil {
		utils.LogErrorCouldNotBindFlag(err)
	}

	if err := viper.BindPFlag(PushAssetsGitHubTokenKey, PipelinePushAssetsCmd.PersistentFlags().Lookup(githubTokenFlag)); err != nil {
		utils.LogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(PushAssetsGithubRepoNameKey, PipelinePushAssetsCmd.PersistentFlags().Lookup(githubRepoNameFlag)); err != nil {
		utils.LogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(PushAssetsGithubUserNameKey, PipelinePushAssetsCmd.PersistentFlags().Lookup(githubUserNameFlag)); err != nil {
		utils.LogErrorCouldNotBindFlag(err)
	}

	viper.AutomaticEnv()

	PipelinePushCmd.AddCommand(PipelinePushAssetsCmd)
}
