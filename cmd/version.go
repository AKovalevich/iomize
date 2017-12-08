package cmd

import (
	"fmt"

	"github.com/AKovalevich/iomize/pkg/version"
	"github.com/spf13/cobra"
	"runtime"
)

const (
	VersionTemplate = `Version:      %s
Codename:     %s
Go version:   %s
Built:        %s
OS/Arch:      %s`
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Current version of oimize",
	Long:  `Oimize is a CLI library for compression and optimization of images.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(VersionTemplate+"\n",
			version.Current(),
			version.Codename(),
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
		)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
