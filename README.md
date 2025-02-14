# go-cobrautil

CobraUtil is an enhancement to [github.com/spf13/cobra](https://github.com/spf13/cobra) to handle aspects of building a CLI app that I feel Cobra itself ignores.  


## The Missing Features

Features that Cobra itself ignores are:

1. A complete structure for implemented CLI apps,
2. A standardized approach for commands and flags, 
3. A testing framework for testing CLI app commands.

## Usage

To use CobraUtil in a CLI app project:

1. Create a `.go.mod` file with this content:
   ```
   module github.com/your-org/your-repo

   go <latest_version>
   ```
2. Create a `./cmd/yourcli` directory off your project's repo root.

3. Create a `cmds` subdirectory under `./cmd/yourcli`.

4. Create a `./cmd/yourcli/go.mod` file with this content:
   ```
   module github.com/your-org/your-repo/your-cli

   go <latest_version>

   require (
      github.com/spf13/cobra <latest_version>
   )
   ```
5. Create a `./cmd/yourcli/go.work` file with this content:
   ```
   go 1.24
   
   use (
      .
      ./../..
   )
   ```   
6. Create a `./cmd/yourcli/main.go` file and add the following code:

   ```go
   package main
   
   import (
      "os"
   
      "github.com/mikeschinkel/go-cobrautil"
      "github.com/your-org/your-repo/your-cli/cmds"
   )
   
   func main() {
      _, err := cobrautil.Execute(cmds.RootCmd, cobrautil.ExecuteArgs{})
      if err != nil {
         os.Exit(1)
      }
   }
   
   ```

7. Create a `./yourpkg/flags.go` file and add the following:
  ```
  package yourpkg
  var GlobalFlags struct {
    Quiet  *bool
    Output *string
  }
  ```

1. Create a `./cmd/yourcli/cmd_root.go` file and add the following code:

   ```go
   package cmds
   
   import (
     "github.com/mikeschinkel/go-cobrautil"
     "github.com/spf13/cobra"
   )
   
   // RootCmd represents the base command when called without any subcommands
   var RootCmd = NewCmdFromOpts(CmdOpts{
   Command: &cobra.Command{
     Use:   "prefsctl",
     Short: "CLI for managing macOS preferences",
     Long:  "CLI for managing macOS preferences, especially for use with Ansible",
     PersistentPreRun: func(cmd *cobra.Command, args []string) {
       cobrautil.SetCalledCmd(cmd)
       if *macprefs.GlobalFlags.Quiet {
         cobrautil.SetQuiet(cmd)
       }
     },
     // Silence usage as we present usage after calling Cobra
     SilenceUsage: true,
     // Silence errors as we handle errors after calling Cobra
     SilenceErrors: true,
     },
   })
   
   func init() {
     cobrautil.AddInitializer(func(cli *CLI) {
     
       yourpkg.GlobalFlags.YourGlobalFlag = PersistentBoolFlag(RootCmd, CmdFlagArgs{
         Name:      cobrautil.QuietFlagName,
         Shorthand: cobrautil.QuietFlagShorthand,
         Default:   false,
         Usage:     "Disable informational messages to stdOut",
       })
     
       yourpkg.GlobalFlags.Output = PersistentStringFlag(RootCmd, CmdFlagArgs{
         Name:      cobrautil.OutputFlagName,
         Shorthand: cobrautil.OutputFlagShorthand,
         Default:   "",
         Usage: fmt.Sprintf("Specify the format for output; one of: %s",
           strings.Join(sliceconv.ToStrings(yourpkg.AllFormats), ","),
         ),
       })
     })
   }

   ```
   
9. Create a `./yourpkg` directory off the root of your repo. This directory is for the package you will develop that implements your CLI's functionality and that your commands will call, commands that will be stored in `./cmd/yourcli/cmds`:

### Usage for each command
Then for each command:

10. Create a file in `./yourpkg` named whatever you prefer to store `YourCmdAction()` with the following signature:
   ```
   YourCmdAction(ctx Context, cfg cobrautil.Config, cmd Cmd) cobrautil.CmdResult
   ```
11. Create a file in `./cmd/yourcli/cmds` named in the form of `cmd_yourcmd.go` and add the following code, changing `yourcmd` and `yourCmd` and `yourpkg` to names applicable to your command and the package which provides your command's functionality:   
   ```go
   package cmds
   
   import (
     "reflect"
   
     "github.com/mikeschinkel/go-cobrautil"
     "github.com/your-org/your-repo/yourpkg"
   )
   
   func init() {
     cobrautil.AddInitializer(func(cli *CLI) {
       RootCmd.AddCmd(yourcmdCmd)
     })
   }
   
   type YourCmdProps struct {
     BaseProps
     yourProp yourpkg.YourPropertyType
   }
   
   var yourCmdProps = &YourCmdProps{}
   
   // yourCmd represents the `yourcmd` command
   var yourCmd = NewCmdFromOpts(CmdOpts{
     Parent: RootCmd,
     Command: &Command{
       Use:                       "YourCmd -y <yourflag>",
       Short:                     "A short description goes here",
       Long:                      "A longer description goes here",
       <OtherCobraCommandProps>:  <theirValues>,
     },
     Props: yourcmdProps,
     Flags: []*CmdFlag{
       {
         Name:      "yourflag",
         Type:      reflect.String,    // or whatever type you need
         Required:  true,              // or omit if not required
         Shorthand: 'y',
         AssignFunc: func(value any) {
           yourCmdProps.yourProp = value.(*string) // Or *type for other types
         },
       },
     },
     RunFunc: runYourCmdFunc,
   })
   
   func runYourCmdFunc(ctx Context, cfg cobrautil.Config, cmd Cmd) cobrautil.CmdResult {
     return yourpkg.YourCmdAction(ctx, cfg, cmd, yourpkg.YourCmdArgs{
       YourProp: yourpkg.YourProp(*cmd.Props().(*YourCmdProps).yourProp),
     }).CobraUtilResult(cmd)
   }
   ```

## Testing
Coming soon...