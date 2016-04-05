package main

import (
  "github.com/codegangsta/cli"
  "os"
  "fmt"
  "io/ioutil"
  "os/user"
  "path"
  "regexp"
  "os/exec"
)

func main() {
  app := cli.NewApp()
  app.Name = "Launch"
  app.Usage = "Convenient wrapper around launchctl"
  app.Commands = []cli.Command {
  {
    Name: "list",
    Aliases: []string{"ls"},
    Usage: "List all plist items. Optionally matching PATTERN",
    Flags: []cli.Flag {
      cli.BoolFlag {
        Name: "full, f",
        Usage: "Show full file paths",
      },
    },
    Action: func(c *cli.Context) {
      var pattern string = c.Args().First()
      full := c.Bool("full")

      for short, long := range filter_plists(pattern) {
        if full {
          fmt.Println(long)
        } else {
          fmt.Println(short)
        }
      }
    },
  },
  {
    Name: "start",
    Usage: "Start a plist",
    Action: func(c *cli.Context) {
      execute_command("load", c)
    },
  },
  {
    Name: "stop",
    Usage: "Stop a plist",
    Action: func(c *cli.Context) {
      execute_command("unload", c)
    },
  },
  {
    Name: "restart",
    Usage: "Restart a plist",
    Action: func(c *cli.Context) {
      execute_command("unload", c)
      execute_command("load", c)
    },
  },
}

  app.Run(os.Args)
}

// Executes a command for a single plist item.
func execute_command(command string, c *cli.Context) {
  var pattern string = c.Args().First()
  plist := single_filtered_plist(pattern)

  out, err := exec.Command("launchctl", command, plist).CombinedOutput()

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  fmt.Print(string(out))
}

// Return a single plist file path to execute a command on.
func single_filtered_plist(pattern string) string {
  filtered_plists := filter_plists(pattern)

  if len(filtered_plists) == 0 {
    fmt.Println("No matches")
    os.Exit(1)
  }

  if len(filtered_plists) > 1 {
    fmt.Println("Too many matches")
    os.Exit(1)
  }

  for _, plist := range filter_plists(pattern) {
    return plist;
  }

  return ""
}

// Filter plist items by a pattern.
func filter_plists(pattern string) map[string]string {
  // Return all plists if pattern is empty.
  if pattern == "" {
    return plists()
  }

  filtered_plists := make(map[string]string)
  regexp := regexp.MustCompile(pattern)

  for k, v := range plists() {
    if regexp.MatchString(k) {
      filtered_plists[k] = v
    }
  }

  return filtered_plists
}

// Get a list of plist files in each plist directory.
func plists() map[string]string {
  plists := make(map[string]string)

  for _, dir := range dirs() {
    // Check if the dir exists. If it doesn't, move on to the next.
    if _, err := os.Stat(dir); err != nil {
      continue
    }

    plist_files, err := ioutil.ReadDir(dir)

    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    for _, file := range plist_files {
      plists[file.Name()] = path.Join(dir, file.Name())
    }
  }

  return plists
}

// Get a list of direcotries to search.
func dirs() []string {
  user, err := user.Current()

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  dirs := []string{"/Library/LaunchAgents", path.Join(user.HomeDir, "Library", "LaunchAgents")}

  if is_root() {
    root_dirs := []string{"/Library/LaunchDaemons", "/System/Library/LaunchDaemons"}
    dirs = append(dirs, root_dirs ...)
  }

  return dirs
}

// Determine if the current process is running as root.
func is_root() bool {
  return os.Geteuid() == 0
}
