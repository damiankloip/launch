package main

import (
  "github.com/codegangsta/cli"
  "github.com/fatih/color"
  "os"
  "path"
  "io"
  "fmt"
  "os/exec"
)

// Executes a command for a single plist item.
func execute_command(command string, c *cli.Context) {
  pattern := c.Args().First()
  plist := single_filtered_plist(pattern)

  write := c.Bool("write")

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

  if (len(out) > 0) {
    yellow := color.New(color.FgYellow)
    yellow.Print("  ", string(out))
  }
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
      fmt.Printf("%s installed to %s (linked)\n", file, dest_path)
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

func uninstall(c *cli.Context) {
  pattern := c.Args().First()
  plist := single_filtered_plist(pattern)

  err := os.Remove(plist)
  check_error(err)

  fmt.Printf("Removed %s\n", plist)
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
