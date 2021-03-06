package cmd

import (
	"log"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"

	"github.com/Shopify/themekit/kit"
)

var (
	arbiter          = newCommandArbiter()
	bootstrapVersion string
	bootstrapPrefix  string
	bootstrapURL     string
	bootstrapName    string
	setThemeID       bool
	openEdit         bool
	openWith         string
	updateVersion    string
	noUpdateNotifier bool

	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()

	stdOut = log.New(colorable.NewColorableStdout(), "", 0)
	stdErr = log.New(colorable.NewColorableStderr(), "", 0)
)

// ThemeCmd is the main entry point to the theme kit command line interface.
var ThemeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Theme Kit is a tool kit for manipulating shopify themes",
	Long: `Theme Kit is a tool kit for manipulating shopify themes

Theme Kit is a fast and cross platform tool that enables you
to build shopify themes with ease.

Complete documentation is available at https://shopify.github.io/themekit/`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !noUpdateNotifier && kit.IsNewUpdateAvailable() {
			stdOut.Print(yellow("An update for Themekit is available. To update please run `theme update`"))
		}
		arbiter.setFlagConfig()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		arbiter.progress.Stop()
	},
}

func init() {
	ThemeCmd.PersistentFlags().StringVarP(&arbiter.configPath, "config", "c", arbiter.configPath, "path to config.yml")
	ThemeCmd.PersistentFlags().VarP(&arbiter.environments, "env", "e", "environment to run the command")
	ThemeCmd.PersistentFlags().StringVarP(&arbiter.flagConfig.Directory, "dir", "d", "", "directory that command will take effect. (default current directory)")
	ThemeCmd.PersistentFlags().StringVarP(&arbiter.flagConfig.Password, "password", "p", "", "theme password. This will override what is in your config.yml")
	ThemeCmd.PersistentFlags().StringVarP(&arbiter.flagConfig.ThemeID, "themeid", "t", "", "theme id. This will override what is in your config.yml")
	ThemeCmd.PersistentFlags().StringVarP(&arbiter.flagConfig.Domain, "store", "s", "", "your shopify domain. This will override what is in your config.yml")
	ThemeCmd.PersistentFlags().StringVar(&arbiter.flagConfig.Proxy, "proxy", "", "proxy for all theme requests. This will override what is in your config.yml")
	ThemeCmd.PersistentFlags().DurationVar(&arbiter.flagConfig.Timeout, "timeout", 0, "the timeout to kill any stalled processes. This will override what is in your config.yml")
	ThemeCmd.PersistentFlags().BoolVarP(&arbiter.verbose, "verbose", "v", false, "Enable more verbose output from the running command.")
	ThemeCmd.PersistentFlags().BoolVarP(&noUpdateNotifier, "no-update-notifier", "", false, "Stop theme kit from notifying about updates.")
	ThemeCmd.PersistentFlags().Var(&arbiter.ignoredFiles, "ignored-file", "A single file to ignore, use the flag multiple times to add multiple.")
	ThemeCmd.PersistentFlags().Var(&arbiter.ignores, "ignores", "A path to a file that contains ignore patterns.")
	ThemeCmd.PersistentFlags().BoolVar(&arbiter.disableIgnore, "no-ignore", false, "Will disable config ignores so that all files can be changed")

	watchCmd.Flags().StringVarP(&arbiter.notifyFile, "notify", "n", "", "file to touch when workers have gone idle")
	watchCmd.Flags().BoolVarP(&arbiter.allenvs, "allenvs", "a", false, "run command with all environments")
	removeCmd.Flags().BoolVarP(&arbiter.allenvs, "allenvs", "a", false, "run command with all environments")
	replaceCmd.Flags().BoolVarP(&arbiter.allenvs, "allenvs", "a", false, "run command with all environments")
	uploadCmd.Flags().BoolVarP(&arbiter.allenvs, "allenvs", "a", false, "run command with all environments")
	openCmd.Flags().BoolVarP(&arbiter.allenvs, "allenvs", "a", false, "run command with all environments")

	downloadCmd.Flags().BoolVarP(&arbiter.force, "force", "f", false, "disable version checking and force all changes")
	watchCmd.Flags().BoolVarP(&arbiter.force, "force", "f", false, "disable version checking and force all changes")
	removeCmd.Flags().BoolVarP(&arbiter.force, "force", "f", false, "disable version checking and force all changes")
	replaceCmd.Flags().BoolVarP(&arbiter.force, "force", "f", false, "disable version checking and force all changes")
	uploadCmd.Flags().BoolVarP(&arbiter.force, "force", "f", false, "disable version checking and force all changes")

	watchCmd.Flags().StringVarP(&arbiter.master, "master", "m", kit.DefaultEnvironment, "The destination from which all changes will be applied")
	removeCmd.Flags().StringVarP(&arbiter.master, "master", "m", kit.DefaultEnvironment, "The destination from which all changes will be applied")
	replaceCmd.Flags().StringVarP(&arbiter.master, "master", "m", kit.DefaultEnvironment, "The destination from which all changes will be applied")
	uploadCmd.Flags().StringVarP(&arbiter.master, "master", "m", kit.DefaultEnvironment, "The destination from which all changes will be applied")

	bootstrapCmd.Flags().StringVar(&bootstrapVersion, "version", latestRelease, "version of Shopify Timber to use")
	bootstrapCmd.Flags().StringVar(&bootstrapPrefix, "prefix", "", "prefix to the Timber theme being created")
	bootstrapCmd.Flags().StringVar(&bootstrapURL, "url", "", "a url to pull a project theme zip file from.")
	bootstrapCmd.Flags().StringVar(&bootstrapName, "name", "", "a name to define your theme on your shopify admin")

	updateCmd.Flags().StringVar(&updateVersion, "version", "latest", "version of themekit to install")

	openCmd.Flags().BoolVarP(&openEdit, "edit", "E", false, "open the web editor for the theme.")
	openCmd.Flags().StringVarP(&openWith, "browser", "b", "", "name of the browser to open the url. the name should match the name of browser on your system.")

	ThemeCmd.AddCommand(
		bootstrapCmd,
		configureCmd,
		removeCmd,
		replaceCmd,
		uploadCmd,
		watchCmd,
		downloadCmd,
		versionCmd,
		updateCmd,
		openCmd,
	)
}
