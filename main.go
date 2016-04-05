package main

import (
  "github.com/codegangsta/cli"
  "os"
  "fmt"
  "io/ioutil"
  "os/user"
  "path"
  "regexp"
)

func main() {
  app := cli.NewApp()
  app.Name = "Launch"
  app.Usage = "Convenient wrapper around launchctl"
  app.Commands = []cli.Command{
  {
    Name: "list",
    Aliases: []string{"ls"},
    Usage: "add a task to the list",
    Action: func(c *cli.Context) {
      var pattern string = c.Args().First()

      for _, plist := range filter_plists(pattern) {
        fmt.Println(plist)
      }
    },
  },
  {
    Name: "start",
    Usage: "complete a task on the list",
    Action: func(c *cli.Context) {
      println("completed task: ", c.Args().First())
    },
  },
  {
    Name: "stop",
    Usage: "complete a task on the list",
    Action: func(c *cli.Context) {
      println("completed task: ", c.Args().First())
    },
  },
  {
    Name: "restart",
    Usage: "complete a task on the list",
    Action: func(c *cli.Context) {
      println("completed task: ", c.Args().First())
    },
  },
}

  app.Run(os.Args)
}

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
