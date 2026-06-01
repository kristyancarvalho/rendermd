package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kristyancarvalho/mdp/internal/config"
	"github.com/kristyancarvalho/mdp/internal/ui"
	"github.com/kristyancarvalho/mdp/internal/watch"
)

func Run(version string, args []string) error {
	fs := flag.NewFlagSet("mdp", flag.ContinueOnError)
	var (
		watchFlag  = fs.Bool("watch", false, "enable hot reload on file change")
		configFlag = fs.String("config", config.DefaultPath(), "path to config file")
		themeFlag  = fs.String("theme", "", "theme name or path")
		versionFlag = fs.Bool("version", false, "print version and exit")
	)
	fs.BoolVar(watchFlag, "w", false, "enable hot reload (short)")
	fs.StringVar(configFlag, "c", config.DefaultPath(), "path to config (short)")
	fs.StringVar(themeFlag, "t", "", "theme (short)")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	if *versionFlag {
		fmt.Println("mdp", version)
		return nil
	}

	cfg := config.Load(*configFlag)
	if *themeFlag != "" {
		cfg.Theme.Name = *themeFlag
	}

	var (
		filename string
		content  []byte
		err      error
	)

	posArgs := fs.Args()
	if len(posArgs) > 0 {
		filename = posArgs[0]
		content, err = os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("cannot read file: %w", err)
		}
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			content, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("cannot read stdin: %w", err)
			}
			*watchFlag = false
		} else {
			fmt.Fprintf(os.Stderr, "usage: mdp [--watch] [--config <path>] [--theme <name>] <file>\n")
			os.Exit(1)
		}
	}

	var w *watch.Watcher
	if *watchFlag && filename != "" {
		debounce := time.Duration(cfg.Watch.DebounceMs) * time.Millisecond
		w, err = watch.New(filename, debounce)
		if err != nil {
			return fmt.Errorf("cannot watch file: %w", err)
		}
		defer w.Close()
	}

	return ui.Run(filename, content, cfg, w)
}
