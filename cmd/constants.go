package cmd

const (
	LangGo      = "go"
	LangDefault = LangGo

	SyncClientRedisUrlDefault      = "localhost:6379"
	SyncClientRedisPrefixDefault   = "dibs"
	SyncClientRedisPasswordDefault = ""

	SyncClientGitUpCommitMessageUpSynced = "up_synced"
	SyncClientGitUpRemoteName            = "dibs-sync"
	SyncClientGitUpUserName              = "dibs-syncer"
	SyncClientGitUpUserEmail             = "dibs-syncer@pojtinger.space"

	SyncClientPlatformPlaceholder    = "infer"
	SyncClientIgnoreRegexPlaceholder = "infer"

	SyncRedisSuffixUpBuilt        = "up_built"
	SyncRedisSuffixUpTested       = "up_tested"
	SyncRedisSuffixUpStarted      = "up_started"
	SyncRedisSuffixUpRegistered   = "up_registered"
	SyncRedisSuffixUpUnregistered = "up_unregistered"
	SyncRedisSuffixUpPushed       = "up_pushed"

	TestIntegrationChartKubernetesIpDefault = "127.0.0.1"

	PlatformAll     = "all"
	PlatformDefault = PlatformAll

	ExecutorNative  = "native"
	ExecutorDocker  = "docker"
	ExecutorDefault = ExecutorNative

	DibsPath        = "."
	DibsName        = ".dibs"
	DibsFileDefault = DibsName + ".yml"

	EnvPrefix = "dibs"

	PlatformEnvDocker = "TARGETPLATFORM" // This is the env variable convention that Docker uses, so alias it
)

// The keys to use for the viper configuration
const (
	LangKey = "lang"

	SyncKeyPrefix = "sync_"

	SyncRedisUrlKey      = SyncKeyPrefix + "redis_url"
	SyncRedisPrefixKey   = SyncKeyPrefix + "redis_prefix"
	SyncRedisPasswordKey = SyncKeyPrefix + "redis_password"

	SyncClientGoGitBaseUrlKey = SyncKeyPrefix + LangGo + "_git_base_url"

	SyncClientPipelineUpDirSrcKey   = SyncKeyPrefix + "dir_src"
	SyncClientPipelineUpDirPushKey  = SyncKeyPrefix + "dir_push"
	SyncClientPipelineUpDirWatchKey = SyncKeyPrefix + "dir_watch"

	SyncClientGoPipelineUpFileModKey = SyncKeyPrefix + LangGo + "_modules_file"

	SyncClientPipelineUpBuildCommandKey = SyncKeyPrefix + "cmd_build"
	SyncClientPipelineUpTestCommandKey  = SyncKeyPrefix + "cmd_test"
	SyncClientPipelineUpStartCommandKey = SyncKeyPrefix + "cmd_start"

	SyncClientPipelineUpRegexIgnoreKey    = SyncKeyPrefix + "regex_ignore"
	SyncClientGoPipelineDownModulesKey    = SyncKeyPrefix + LangGo + "_modules_pull"
	SyncClientGoPipelineDownDirModulesKey = SyncKeyPrefix + LangGo + "_dir_pull"

	SyncServerGitServerReposDirKey = SyncKeyPrefix + LangGo + "_dir_repos"
	SyncServerGitServerHttpPortKey = SyncKeyPrefix + LangGo + "-port"
	SyncServerGitServerHttpPathKey = SyncKeyPrefix + LangGo + "-path"

	TestIntegrationChartKubernetesIpKey = "kubernetes_ip"

	PlatformKey = "platform"
	ExecutorKey = "executor"

	DibsFileKey = "config_file"

	PushAssetsKeyPrefix = "push_assets_"

	PushAssetsVersionKey        = PushAssetsKeyPrefix + "version"
	PushAssetsGitHubTokenKey    = PushAssetsKeyPrefix + "github_token"
	PushAssetsGithubUserNameKey = PushAssetsKeyPrefix + "github_user_name"
	PushAssetsGithubRepoNameKey = PushAssetsKeyPrefix + "github_repo_name"

	PushChartKeyPrefix = "push_chart_"

	PushChartGitUserNameKey      = PushChartKeyPrefix + "git_user_name"
	PushChartGitUserEmailKey     = PushChartKeyPrefix + "git_user_email"
	PushChartGitCommitMessageKey = PushChartKeyPrefix + "git_commit_message"
	PushChartGitRepoURLKey       = PushChartKeyPrefix + "git_repo_url"

	PushChartGithubUserNameKey = PushChartKeyPrefix + "github_user_name"
	PushChartGithubTokenKey    = PushChartKeyPrefix + "github_token"
	PushChartGithubRepoNameKey = PushChartKeyPrefix + "github_repo_name"
	PushChartGithubPagesURLKey = PushChartKeyPrefix + "github_pages_url"
)
