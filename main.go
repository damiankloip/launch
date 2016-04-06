package main

import (
  "github.com/codegangsta/cli"
  "os"
  "fmt"
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
