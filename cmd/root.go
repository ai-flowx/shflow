package cmd

import (
	"context"
	"fmt"
	"github.com/ai-flowx/flowx/memory"
	"github.com/ai-flowx/flowx/store"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/ai-flowx/flowx/config"
	"github.com/ai-flowx/flowx/flow"
)

const (
	routineNum = -1
)

var (
	configFile string
	listenAddr string
)

var rootCmd = &cobra.Command{
	Use:     "flowx",
	Version: config.Version + "-build-" + config.Build,
	Short:   "ai framework",
	Long:    fmt.Sprintf("ai framework %s (%s %s)", config.Version, config.Commit, config.Build),
	Run: func(cmd *cobra.Command, args []string) {
		var cfg config.Config
		ctx := context.Background()
		if err := viper.Unmarshal(&cfg); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		s, err := initStore(ctx, &cfg)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		m, err := initMemory(ctx, &cfg, s)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		f, err := initFlow(ctx, &cfg, m)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		if err := runFlow(ctx, f); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// nolint:gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&listenAddr, "listen-addr", "u", "127.0.0.1:8080", "listen address")

	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	if configFile == "" {
		return
	}

	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}

func initStore(ctx context.Context, cfg *config.Config) (store.Store, error) {
	c := store.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config\n")
	}

	c.Provider = cfg.Store.Provider
	c.Api = cfg.Store.Api
	c.Token = cfg.Store.Token

	return store.New(ctx, c), nil
}

func initMemory(ctx context.Context, cfg *config.Config, st store.Store) (memory.Memory, error) {
	c := memory.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config\n")
	}

	c.Store = st

	return memory.New(ctx, c), nil
}

func initFlow(ctx context.Context, _ *config.Config, mem memory.Memory) (flow.Flow, error) {
	c := flow.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config\n")
	}

	c.Addr = listenAddr
	c.Memory = mem

	return flow.New(ctx, c), nil
}

func runFlow(ctx context.Context, fl flow.Flow) error {
	if err := fl.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init\n")
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(routineNum)

	g.Go(func() error {
		if err := fl.Run(ctx); err != nil {
			return errors.Wrap(err, "failed to run\n")
		}
		return nil
	})

	s := make(chan os.Signal, 1)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can"t be caught, so don't need add it
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		<-s
		if err := fl.Deinit(ctx); err != nil {
			return errors.Wrap(err, "failed to deinit\n")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
