package main

import (
  "github.com/codegangsta/cli"
  "github.com/fatih/color"
  "os"
  "path"
  "io"
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
        Usage: "@todo",
      },
    },
    Action: func(c *cli.Context) {
      install(c)
    },
  },
}

  app.Run(os.Args)
}

// Executes a command for a single plist item.
func execute_command(command string, c *cli.Context) {
  var pattern string = c.Args().First()

  if pattern == "" {
    print_error("No plist pattern provided")
  }

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

func install(c *cli.Context) {
  var file string = c.Args().First()

  if file == "" {
    print_error("No file name provided")
  }

  // Check the file exists.
  _, err := os.Stat(file)
  check_error(err)

  if path.Ext(file) != ".plist" {
    print_error("Only files with the .plist extension can be installed")
  }

  symlink := c.Bool("symlink")

  for _, dir := range user_dirs() {
    // Move onto the next if the dir doesn't exist.
    if _, err := os.Stat(dir); err != nil {
      continue
    }

    // Join the source file name with the plist dir.
    dest_path := path.Join(dir, path.Base(file))

    if symlink {
      err := os.Symlink(file, dest_path)
      check_error(err)
      // If we get here, no errors.
      fmt.Print("%s installed to %s\n (linked)", file, dest_path)
    } else {
      // Copy the file to the directory with the original name.
      // Open the source file.
      source, err := os.Open(file)
      check_error(err)

      dest, err := os.Create(dest_path)
      defer dest.Close()
      check_error(err)

      _, err = io.Copy(source, dest)
      check_error(err)

      fmt.Printf("%s installed to %s\n", file, dest_path)
    }

    // Always break if the dir exists. Above means it either successfully
    // wrote/linked a new plist file or failed with an error first.
    break
  }
}

// Checks and handles errors.
func check_error(err error) {
  if err != nil {
    print_error(err)
  }
}

func print_error(message interface{}) {
  fmt.Println(message)
  os.Exit(1)
}
