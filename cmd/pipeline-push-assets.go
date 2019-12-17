package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
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
			if err := platform.Assets.Push(platform.Platform, strings.Split(viper.GetString(AssetsVersionKey), " "), viper.GetString(AssetsGitHubTokenKey)); err != nil {
				utils.PipeLogErrorFatal("Could not push assets", err, platform.Platform)
			}
		}
	},
}

func init() {
	var (
		version string
		token   string

		versionFlag     = strings.Replace(AssetsVersionKey, "assets_", "", -1)
		githubTokenFlag = strings.Replace(AssetsGitHubTokenKey, "_", "-", -1)
	)

	PipelinePushAssetsCmd.PersistentFlags().StringVarP(&version, versionFlag, "v", "0.0.1", `The version of the asset to deploy (use "-prerelease <version>" as the version to create a prerelease)`)
	PipelinePushAssetsCmd.PersistentFlags().StringVarP(&token, githubTokenFlag, "t", "1234", "GitHub personal access token")

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(AssetsVersionKey, PipelinePushAssetsCmd.PersistentFlags().Lookup(versionFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(AssetsGitHubTokenKey, PipelinePushAssetsCmd.PersistentFlags().Lookup(githubTokenFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	viper.AutomaticEnv()

	PipelinePushCmd.AddCommand(PipelinePushAssetsCmd)
}
