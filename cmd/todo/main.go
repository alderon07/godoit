package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// Version info (optional: injected at build time via -ldflags "-X main.Version=1.0.0")
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func usage() {
	fmt.Printf(`godo — minimal todo CLI
Version: %s (built: %s)

Usage:
  godo <command> [options]

Commands:
  add       Add a new task
  list      List tasks
  done      Mark task as complete
  edit      Edit an existing task
  rm        Remove a task
  alerts    Show due/overdue tasks
  stats     Show task analytics
  server    Start HTTP API server
  help      Show this help
  version   Show version info

Run "godo <command> -h" for detailed help on each command.
`, Version, BuildTime)
}

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "add":
		addFlags := flag.NewFlagSet("add", flag.ExitOnError)
		title := addFlags.String("title", "", "Task title (required)")
		description := addFlags.String("desc", "", "Task description (optional)")
		dueStr := addFlags.String("due", "", "Due date YYYY-MM-DD")
		repeat := addFlags.String("repeat", "", "Repeat rule: daily|weekly|monthly")
		priority := addFlags.Int("p", 1, "Priority (1-3, default 1)")
		tags := addFlags.String("tags", "", "Comma-separated tags")
		after := addFlags.String("after", "", "Comma-separated dependency task IDs")
		_ = addFlags.Parse(args)

		RunAdd(*title, *description, *dueStr, *repeat, *priority, *tags, *after)

	case "list", "ls":
		lsFlags := flag.NewFlagSet("list", flag.ExitOnError)
		showAll := lsFlags.Bool("all", false, "Show completed tasks too")
		today := lsFlags.Bool("today", false, "Show only today's tasks")
		week := lsFlags.Bool("week", false, "Show only this week's tasks")
		detailed := lsFlags.Bool("detailed", false, "Show detailed task information")
		grep := lsFlags.String("grep", "", "Filter by substring (case-insensitive)")
		tags := lsFlags.String("tags", "", "Filter by tags (comma=OR, plus=AND)")
		sortKey := lsFlags.String("sort", "due", "Sort by: due|priority|created|status|title")
		before := lsFlags.String("before", "", "Filter tasks before YYYY-MM-DD")
		after := lsFlags.String("after", "", "Filter tasks after YYYY-MM-DD")
		_ = lsFlags.Parse(args)

		RunList(*showAll, *today, *week, *detailed, *grep, *tags, *sortKey, *before, *after)

	case "done":
		doneFlags := flag.NewFlagSet("done", flag.ExitOnError)
		_ = doneFlags.Parse(args)

		if doneFlags.NArg() < 1 {
			log.Fatal("Usage: godo done <index>")
		}

		RunDone(doneFlags.Arg(0))

	case "edit":
		editFlags := flag.NewFlagSet("edit", flag.ExitOnError)
		title := editFlags.String("title", "", "New task title")
		description := editFlags.String("desc", "", "New task description (or 'none' to clear)")
		dueStr := editFlags.String("due", "", "Due date YYYY-MM-DD (or 'none' to clear)")
		repeat := editFlags.String("repeat", "", "Repeat rule (or 'none' to clear)")
		priority := editFlags.Int("p", 0, "Priority (1-3, 0 to keep current)")
		tags := editFlags.String("tags", "", "Tags (or 'none' to clear)")
		after := editFlags.String("after", "", "Dependencies (or 'none' to clear)")
		_ = editFlags.Parse(args)

		if editFlags.NArg() < 1 {
			log.Fatal("Usage: godo edit <index> [options]")
		}

		RunEdit(editFlags.Arg(0), *title, *description, *dueStr, *repeat, *priority, *tags, *after)

	case "remove", "rm":
		rmFlags := flag.NewFlagSet("rm", flag.ExitOnError)
		_ = rmFlags.Parse(args)

		if rmFlags.NArg() < 1 {
			log.Fatal("Usage: godo rm <index>")
		}

		RunRemove(rmFlags.Arg(0))

	case "alerts":
		alertFlags := flag.NewFlagSet("alerts", flag.ExitOnError)
		watch := alertFlags.Bool("watch", false, "Continuously monitor for upcoming tasks")
		interval := alertFlags.Duration("interval", 60*time.Second, "Polling interval for watch mode")
		ahead := alertFlags.Duration("ahead", 24*time.Hour, "Lookahead window for alerts")
		_ = alertFlags.Parse(args)

		RunAlerts(*watch, *interval, *ahead)

	case "stats":
		RunStats()

	case "server":
		serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
		port := serverFlags.Int("port", 8080, "Port to listen on")
		host := serverFlags.String("host", "localhost", "Host to bind to")
		_ = serverFlags.Parse(args)

		RunServer(*host, *port)

	case "help", "-h", "--help":
		usage()

	case "version", "-v", "--version":
		fmt.Printf("godo version %s\n", Version)
		fmt.Printf("Build time: %s\n", BuildTime)

	default:
		fmt.Printf("Unknown command: %s\n\n", cmd)
		usage()
		os.Exit(2)
	}
}
