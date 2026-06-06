package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kristyancarvalho/rendermd/internal/config"
	"github.com/kristyancarvalho/rendermd/internal/ui"
	"github.com/kristyancarvalho/rendermd/internal/watch"
)

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

func (b BuildInfo) Print(w io.Writer) {
	fmt.Fprintf(w, "rendermd %s\ncommit %s\nbuilt %s\n", valueOrDefault(b.Version, "dev"), valueOrDefault(b.Commit, "unknown"), valueOrDefault(b.Date, "unknown"))
}

func valueOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func isMarkdownFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".md" || ext == ".markdown"
}

func Run(build BuildInfo, args []string) error {
	fs := flag.NewFlagSet("rendermd", flag.ContinueOnError)
	var (
		watchFlag   = fs.Bool("watch", false, "enable hot reload on file change")
		configFlag  = fs.String("config", config.DefaultPath(), "path to config file")
		themeFlag   = fs.String("theme", "", "theme name or path")
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
		build.Print(os.Stdout)
		return nil
	}

	cfg := config.Load(*configFlag)
	if *themeFlag != "" {
		for _, warning := range cfg.SetThemeName(*themeFlag) {
			fmt.Fprintln(os.Stderr, warning.String())
		}
	}

	var (
		filename string
		content  []byte
		err      error
	)

	posArgs := fs.Args()
	if len(posArgs) > 0 {
		filename = posArgs[0]
		if !isMarkdownFile(filename) {
			return fmt.Errorf("unsupported file type %q: rendermd only accepts .md and .markdown files", filepath.Ext(filename))
		}
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
			fmt.Fprintf(os.Stderr, "usage: rendermd [--watch] [--config <path>] [--theme <name>] <file>\n")
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