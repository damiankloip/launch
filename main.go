package main

import (
  "github.com/codegangsta/cli"
  "github.com/fatih/color"
  "os"
  "fmt"
  "os/exec"
  "io/ioutil"
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
        fmt.Println("No $EDITOR environment variable found.")
        os.Exit(1)
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
}

  app.Run(os.Args)
}

// Executes a command for a single plist item.
func execute_command(command string, c *cli.Context) {
  var pattern string = c.Args().First()
  write := c.Bool("write")
  plist := single_filtered_plist(pattern)

  // Create a slice of command args.
  command_args := []string{command}

  // Add the write flag after the command if needed.
  if (write) {
    command_args = append(command_args, "-w")
  }

  // Add the plist file path.
  command_args = append(command_args, plist)

  command_obj := exec.Command("launchctl", command_args ...)

  fmt.Print("Executing: ")
  green := color.New(color.FgGreen)
  green.Println(command_obj.Args)

  out, err := command_obj.CombinedOutput()
  check_error(err)

  yellow := color.New(color.FgYellow)
  yellow.Print("  ", string(out))
}

// Checks and handles errors.
func check_error(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
