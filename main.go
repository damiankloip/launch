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
    Name: "ls",
    Aliases: []string{"a"},
    Usage: "add a task to the list",
    Action: func(c *cli.Context) {
      for _, plist := range plists() {
        fmt.Println(plist)
      }
    },
  },
  {
    Name: "start",
    Aliases: []string{"c"},
    Usage: "complete a task on the list",
    Action: func(c *cli.Context) {
      println("completed task: ", c.Args().First())
    },
  },
  {
    Name: "stop",
    Aliases: []string{"c"},
    Usage: "complete a task on the list",
    Action: func(c *cli.Context) {
      println("completed task: ", c.Args().First())
    },
  },
  {
    Name: "restart",
    Aliases: []string{"c"},
    Usage: "complete a task on the list",
    Action: func(c *cli.Context) {
      println("completed task: ", c.Args().First())
    },
  },
}

  app.Run(os.Args)
}

func plists() []string {
  var plists []string

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
      plists = append(plists, file.Name())
    }
  }

  return plists
}

//
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
