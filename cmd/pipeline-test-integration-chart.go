package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
	"strings"
)

var PipelineTestIntegrationChartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Integration test the chart",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)
		viperIP := viper.GetString(TestIntegrationChartKubernetesIpKey)

		rawIP := net.ParseIP(viperIP)
		if rawIP == nil {
			utils.PipeLogErrorFatalCouldNotParseIP(viperIP)
			return
		}
		ip := rawIP.String()

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if viper.GetString(ExecutorKey) == ExecutorDocker {
				if output, err := platform.Tests.Integration.Chart.BuildImage(platform.Platform); err != nil {
					utils.PipeLogErrorFatal("Could not build chart integration test chart", err, platform.Platform, output)
				}
				output, err := platform.Tests.Integration.Chart.StartImage(platform.Platform, struct {
					Key   string
					Value string
				}{
					Key:   "IP",
					Value: ip,
				})
				utils.PipeLogErrorInfo("Chart integration test ran in Docker", err, platform.Platform, output)
			} else {
				output, err := platform.Tests.Integration.Chart.Start(platform.Platform)
				utils.PipeLogErrorInfo("Chart integration test ran", err, platform.Platform, output)
			}
		}
	},
}

func init() {
	var (
		kubernetesIp string

		kubernetesIpFlag = strings.Replace(TestIntegrationChartKubernetesIpKey, "_", "-", -1)
	)

	PipelineTestIntegrationChartCmd.PersistentFlags().StringVarP(&kubernetesIp, kubernetesIpFlag, "i", TestIntegrationChartKubernetesIpDefault, "IP of the Kubernetes cluster to create if running in Docker (often the host machine's IP)")

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(TestIntegrationChartKubernetesIpKey, PipelineTestIntegrationChartCmd.PersistentFlags().Lookup(kubernetesIpFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}

	viper.AutomaticEnv()

	PipelineTestIntegrationCmd.AddCommand(PipelineTestIntegrationChartCmd)
}
