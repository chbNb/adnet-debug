package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.mobvista.com/ADN/adnet/cmd/hb_server/commands"
)

var (
	// Version is the current version of Tendermint
	// Must be a string because scripts like dist.sh read this file.
	Version = "n/a"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit = "n/a"

	// GitTag is the current Tag
	GitTag = "0"

	// BuildTime is the binaly build time
	BuildTime = "n/a"
)

func init() {
	Version = fmt.Sprintf("%s built on %s (commit: %s)", GitTag, BuildTime, GitCommit)
	commands.Appversion = GitTag
}

// VersionCmd VersionCmd
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func main() {
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(commands.NewServe())
	rootCmd.AddCommand(VersionCmd)
	rootCmd.PersistentFlags().StringP("test.coverprofile", "", "", "code cover file path")
	rootCmd.PersistentFlags().BoolP("systemTest", "", false, "test flag")
	viper.BindPFlags(rootCmd.Flags())
	rootCmd.Execute()
}
