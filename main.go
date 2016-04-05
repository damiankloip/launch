package main

import (
  "github.com/codegangsta/cli"
  "os"
  "fmt"
  "io/ioutil"
  "os/user"
  "path"
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
      for _, plist := range plists() {
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

func is_root() bool {
  return os.Geteuid() == 0
}
