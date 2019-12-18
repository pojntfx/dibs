package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var PipelinePushAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Push assets",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if output, err := platform.Assets.Push(platform.Platform, strings.Split(viper.GetString(PushAssetsVersionKey), " "), viper.GetString(PushAssetsGitHubTokenKey)); err != nil {
				utils.PipeLogErrorFatal("Could not push assets", err, platform.Platform, output)
			}
		}
	},
}

func init() {
	var (
		version string
		token   string

		versionFlag     = strings.Replace(strings.Replace(PushAssetsVersionKey, PushAssetsKeyPrefix, "", -1), "assets_", "", -1)
		githubTokenFlag = strings.Replace(strings.Replace(PushAssetsGitHubTokenKey, PushAssetsKeyPrefix, "", -1), "_", "-", -1)
	)

	PipelinePushAssetsCmd.PersistentFlags().StringVarP(&version, versionFlag, "v", "0.0.1", `The version of the asset to deploy (use "-prerelease <version>" as the version to create a prerelease)`)
	PipelinePushAssetsCmd.PersistentFlags().StringVarP(&token, githubTokenFlag, "t", "1234", "GitHub personal access token")

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(PushAssetsVersionKey, PipelinePushAssetsCmd.PersistentFlags().Lookup(versionFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(PushAssetsGitHubTokenKey, PipelinePushAssetsCmd.PersistentFlags().Lookup(githubTokenFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}

	viper.AutomaticEnv()

	PipelinePushCmd.AddCommand(PipelinePushAssetsCmd)
}
