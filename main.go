package main

import (
  "github.com/codegangsta/cli"
  "fmt"
  "os"
  "io/ioutil"
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
        Name: "long, l",
        Usage: "Show full file paths",
      },
    },
    Action: func(c *cli.Context) {
      var pattern string = c.Args().First()
      full := c.Bool("long")

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
    Flags: []cli.Flag {
      cli.BoolFlag {
        Name: "write, w",
        Usage: "persist the start behaviour so the agent will load on startup",
      },
    },
    Action: func(c *cli.Context) {
      execute_command("load", c)
    },
  },
  {
    Name: "stop",
    Usage: "Stop a plist",
    Flags: []cli.Flag {
      cli.BoolFlag {
        Name: "write, w",
        Usage: "persist the stop behaviour so the agent will never load on startup",
      },
    },
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
  {
    Name: "show",
    Usage: "Show the contents of a plist file",
    Action: func(c *cli.Context) {
      var pattern string = c.Args().First()
      plist := single_filtered_plist(pattern)

      data, err := ioutil.ReadFile(plist)
      check_error(err)

      fmt.Printf("%s", data)
    },
  },
  {
    Name: "edit",
    Usage: "Edit the contents of a plist file",
    Action: func(c *cli.Context) {
      editor := os.Getenv("EDITOR");

      if editor == "" {
        print_error("No $EDITOR environment variable found.")
      }

      var pattern string = c.Args().First()
      plist := single_filtered_plist(pattern)

      command_obj := exec.Command(editor, plist)

      command_obj.Stdin = os.Stdin
      command_obj.Stdout = os.Stdout
      command_obj.Stderr = os.Stderr

      err := command_obj.Start()
      check_error(err)

      err = command_obj.Wait()
      check_error(err)

      fmt.Printf("Editing %s\n", plist)
    },
  },
  {
    Name: "install",
    Usage: "Install a plist to ~/Library/LaunchAgents or /Library/LaunchAgents (whichever it finds first)",
    Flags: []cli.Flag {
      cli.BoolFlag {
        Name: "symlink, s",
        Usage: "Symlink the source plist file to the LauchAgents directory (instead of copying)",
      },
    },
    Action: func(c *cli.Context) {
      install(c)
    },
  },
  {
    Name: "uninstall",
    Usage: "Uninstall a plist from ~/Library/LaunchAgents or /Library/LaunchAgents (whichever it finds first)",
    Action: func(c *cli.Context) {
      uninstall(c)
    },
  },
}

  app.Run(os.Args)
}
