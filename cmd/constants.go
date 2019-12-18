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

const (
	LangKey = "lang"

	SyncRedisUrlKey      = "redis_url"
	SyncRedisPrefixKey   = "redis_prefix"
	SyncRedisPasswordKey = "redis_password"

	SyncClientGoGitBaseUrlKey = LangGo + "_git_base_url"

	SyncClientPipelineUpDirSrcKey   = "dir_src"
	SyncClientPipelineUpDirPushKey  = "dir_push"
	SyncClientPipelineUpDirWatchKey = "dir_watch"

	SyncClientGoPipelineUpFileModKey = LangGo + "_modules_file"

	SyncClientPipelineUpBuildCommandKey = "cmd_build"
	SyncClientPipelineUpTestCommandKey  = "cmd_test"
	SyncClientPipelineUpStartCommandKey = "cmd_start"

	SyncClientPipelineUpRegexIgnoreKey    = "regex_ignore"
	SyncClientGoPipelineDownModulesKey    = LangGo + "_modules_pull"
	SyncClientGoPipelineDownDirModulesKey = LangGo + "_dir_pull"

	SyncServerGitServerReposDirKey = LangGo + "_dir_repos"
	SyncServerGitServerHttpPortKey = LangGo + "-port"
	SyncServerGitServerHttpPathKey = LangGo + "-path"

	TestIntegrationChartKubernetesIpKey = "kubernetes_ip"

	PlatformKey = "platform"
	ExecutorKey = "executor"

	DibsFileKey = "config_file"

	PushAssetsKeyPrefix = "push_assets_"

	PushAssetsVersionKey     = PushAssetsKeyPrefix + "version"
	PushAssetsGitHubTokenKey = PushAssetsKeyPrefix + "github_token"

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
