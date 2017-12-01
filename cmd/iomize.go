package cmd

import (
	"fmt"
	"os"
	"strconv"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/AKovalevich/iomize/pkg/config"
	//"github.com/AKovalevich/iomize/pkg/pipeline"
	"github.com/AKovalevich/iomize/pkg/server"
	log "github.com/AKovalevich/scrabbler/log/logrus"
	"github.com/AKovalevich/iomize/pkg/pipeline"
	"github.com/AKovalevich/iomize/pkg/entrypoint"
)

var mainConfig *config.MainConfiguration

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "serve",
	//TraverseChildren: true,
	Short: "Oimize",
	Long: `Oimizer is a CLI library for compression and optimization image.`,
	Run: func(cmd *cobra.Command, args []string) {
		start(mainConfig)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	mainConfig = config.NewConfiguration()
	RootCmd.PersistentFlags().StringVar(&mainConfig.ConfigFilePath, "config", "", "Config file (default is $HOME/.test.yaml)")
	RootCmd.Flags().StringVarP(&mainConfig.Port, "port", "p", "1122", "HTTP server port")
	RootCmd.Flags().StringVarP(&mainConfig.Host, "host", "", "localhost", "Host of HTTP server")
	RootCmd.Flags().BoolVarP(&mainConfig.Debug, "debug", "d", false, "Enable debug version")
	RootCmd.Flags().IntVarP(&mainConfig.GraceTimeOut, "grace_time_out", "g", 10, "Graceful timeout in milliseconds")
	RootCmd.Flags().StringVarP(&mainConfig.PipeConfigPath, "pipe_config_path", "", "pipes.yaml", "Pipelines configuration file path (default is $HOME/develop.pipeline.yaml)")
	viper.BindPFlag("port", RootCmd.Flags().Lookup("port"))
	viper.BindPFlag("host", RootCmd.Flags().Lookup("host"))
	viper.BindPFlag("debug", RootCmd.Flags().Lookup("debug"))
	viper.BindPFlag("grace_time_out", RootCmd.Flags().Lookup("grace_time_out"))
	viper.BindPFlag("pipe_config_path", RootCmd.Flags().Lookup("pipe_config_path"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("yaml")
	if mainConfig.ConfigFilePath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(mainConfig.ConfigFilePath)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".test" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".test")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Start scrabbler application
func start(config *config.MainConfiguration) {
	numCPU := runtime.NumCPU()
	log.Do.Info("Set max proccess: " + strconv.Itoa(numCPU))
	runtime.GOMAXPROCS(numCPU)

	log.Do.Info("Initialize entrypoints: ")
	// @TODO Hardcoded list of entrypoints. Need to refactor this
	entryPointList := entrypoint.EntrypointList{}
	log.Do.Info("Initialize \"backend\" entry point")
	entryPointList.Set("backend")
	config.EntryPoints = entryPointList

	log.Do.Info("Initialize pipeline lists")
	var err error
	config.PipeLineList, err = pipeline.InitPipelines(config.PipeConfigPath)
	if err != nil {
		log.Do.Fatal(err.Error())
	}
	log.Do.Debugf("PID: %d\n", os.Getpid())
	s := server.NewServer(config)
	s.Serve()
}
