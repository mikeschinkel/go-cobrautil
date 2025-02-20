package cobrautil

// TODO: Extract out into a proper config package
import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	ConfigFileExt = "yaml"
)

const (
	ConfigFilepath    = "config_filepath"
	ConfigOptions     = "config_options"
	ConfigInitialized = "config_initialized"
)

type ConfigArgs struct {
	AppName  string
	Filepath string
	Filename string
	Options  OptionsMap
}

type Config interface {
	Initialize(Context) error
	AppName() string
	Dir() string
	Filepath() string
	OtherFilepath(string) string
	Options() OptionsMap
}

var _ Config = (*config)(nil)

type config struct {
	dir         string
	appName     string
	filepath    string
	options     OptionsMap
	initialized bool
}

func (c *config) OtherFilepath(filename string) string {
	return filepath.Join(c.dir, filename) //TODO implement me
}

func (c *config) Dir() string {
	return c.dir
}

func (c *config) AppName() string {
	return c.appName
}

//goland:noinspection GoUnusedExportedFunction
func ConfigLogArgs(c config) []any {
	return []any{
		ConfigFilepath, c.Filepath(),
		ConfigOptions, c.Options(),
		ConfigInitialized, c.Initialized(),
	}
}

func (c *config) Options() OptionsMap {
	return c.options
}

func (c *config) Filepath() string {
	return c.filepath
}

func NewConfig(args ConfigArgs) Config {
	if args.Options == nil {
		args.Options = make(OptionsMap)
	}
	return &config{
		appName: args.AppName,
		options: args.Options,
	}
}

func (c *config) Initialized() bool {
	return c.initialized
}

func (c *config) Initialize(_ Context) (err error) {
	var execPath string

	if c.initialized {
		goto end
	}
	if c.appName == "" {
		execPath, err = os.Executable()
		if err != nil {
			goto end
		}
		c.appName = filepath.Base(execPath)
	}
	c.dir, err = ConfigDir(c.appName)
	if err != nil {
		goto end
	}
	err = os.MkdirAll(c.dir, os.ModePerm)
	if err != nil {
		goto end
	}
	c.filepath = filepath.Join(
		c.dir,
		fmt.Sprintf("%s.%s", c.appName, ConfigFileExt),
	)
	c.initialized = true
end:
	return nil
}

func SaveConfig(ctx Context, cfg Config, file string) error {
	return marshalJSONFile(ctx, file, cfg)
}

type OnErrFunc func(string, error)

//goland:noinspection GoUnusedExportedFunction
func SaveConfigWithOnErr(ctx Context, cfg Config, onErr OnErrFunc) (err error) {
	fp := cfg.Filepath()
	err = SaveConfig(ctx, cfg, fp)
	if err != nil {
		onErr(fp, err)
	}
	return err
}

func ConfigDir(configName string) (dir string, err error) {
	dir = os.Getenv("XDG_CONFIG_HOME")
	if dir != "" {
		goto end
	}
	dir = os.Getenv("HOME")
	if dir == "" {
		err = ErrNoHomeDirVar
		goto end
	}
	dir = filepath.Join(dir, ".config", configName)
end:
	return dir, err
}
