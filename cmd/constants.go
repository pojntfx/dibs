package cmd

const (
	LangGo      = "go"
	LangDefault = LangGo

	RedisUrlDefault      = "localhost:6379"
	RedisPrefixDefault   = "dibs"
	RedisPasswordDefault = ""

	LangKey = "lang"

	RedisUrlKey      = "redis_url"
	RedisPrefixKey   = "redis_prefix"
	RedisPasswordKey = "redis_password"

	GitUpCommitMessage = "up_synced"
	GitUpRemoteName    = "dibs-sync"
	GitUpUserName      = "dibs-syncer"
	GitUpUserEmail     = "dibs-syncer@pojtinger.space"

	PlatformPlaceholder    = "infer"
	IgnoreRegexPlaceholder = "infer"

	GoGitBaseUrlKey = LangGo + "_git_base_url"

	PipelineUpDirSrcKey   = "dir_src"
	PipelineUpDirPushKey  = "dir_push"
	PipelineUpDirWatchKey = "dir_watch"

	GoPipelineUpFileModKey = LangGo + "_modules_file"

	PipelineUpBuildCommandKey = "cmd_build"
	PipelineUpTestCommandKey  = "cmd_test"
	PipelineUpStartCommandKey = "cmd_start"

	PipelineUpRegexIgnoreKey    = "regex_ignore"
	GoPipelineDownModulesKey    = LangGo + "_modules_pull"
	GoPipelineDownDirModulesKey = LangGo + "_dir_pull"

	RedisSuffixUpBuilt        = "up_built"
	RedisSuffixUpTested       = "up_tested"
	RedisSuffixUpStarted      = "up_started"
	RedisSuffixUpRegistered   = "up_registered"
	RedisSuffixUpUnregistered = "up_unregistered"
	RedisSuffixUpPushed       = "up_pushed"

	GitServerReposDirKey = LangGo + "_dir_repos"
	GitServerHttpPortKey = LangGo + "-port"
	GitServerHttpPathKey = LangGo + "-path"

	KubernetesIpKey     = "kubernetes_ip"
	KubernetesIpDefault = "127.0.0.1"

	PlatformAll     = "all"
	PlatformDefault = PlatformAll

	ExecutorNative  = "native"
	ExecutorDocker  = "docker"
	ExecutorDefault = ExecutorNative

	DibsPath        = "."
	DibsName        = ".dibs"
	DibsFileDefault = DibsName + ".yml"

	EnvPrefix = "dibs"

	PlatformKey = "platform"
	ExecutorKey = "executor"

	PlatformEnvDocker = "TARGETPLATFORM" // This is the env variable convention that Docker uses, so alias it

	DibsFileKey = "config_file"

	AssetsVersionKey = "assets_version"
	AssetsTokenKey   = "assets_token"
)
