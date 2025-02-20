package cobrautil

import (
	"context"
	"regexp"
	"sync"

	"github.com/spf13/cobra"
)

var rootCmd Cmd
var rootCmdMutex sync.Mutex

var calledCmd *cobra.Command
var calledCmdMutex sync.Mutex

func RootCmd() Cmd {
	return rootCmd
}

//func SetArgsPassed(args []string) {
//	argsPassed = args
//}

type CLI struct {
	Config Config
	Args   []string
}

func NewCLI() *CLI {
	return &CLI{}
}

//goland:noinspection GoUnusedParameter
func (cli *CLI) Execute(ctx Context, args []string) (result CmdResult, err error) {
	rootCmd.SetArgs(args)
	err = rootCmd.Command().Execute()
	result = NewErrorResult(
		NewCmd(cli, calledCmd),
		err,
	)
	return result, err
}

func DefaultContext() Context {
	return context.Background()
}

func (cli *CLI) Initialize(ctx Context, cfg Config, args []string) error {
	var exists bool
	var filepath string

	cli.Config = cfg
	cli.Args = args

	err := ensureDir(cfg.Dir())
	if err != nil {
		goto end
	}
	filepath = cfg.Filepath()

	exists, err = fileExists(filepath)
	if err != nil {
		goto end
	}
	if !exists {
		err = SaveConfig(ctx, cli.Config, filepath)
	}
	if err != nil {
		goto end
	}
	RunInitializers(cli)
end:
	return err
}

var regex = regexp.MustCompile("\n([a-z0-9_]+=)")

func (cli *CLI) ShowUsage(r CmdResult) {
	rootCmd.PrintErr("\nERROR: ")
	rootCmd.PrintErrln(regex.ReplaceAllString(r.Err().Error(), "; $1"))
	rootCmd.PrintErrln("")
	rootCmd.PrintErrln(r.Command().UsageString())
	rootCmd.PrintErrln("")
}
