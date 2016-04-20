package main

import (
  "os"
  "fmt"
  "io/ioutil"
  "os/user"
  "path"
  "regexp"
)

// Return a single plist file path to execute a command on.
func single_filtered_plist(pattern string) string {
  if pattern == "" {
    print_error("No plist pattern provided")
  }

  filtered_plists := filter_plists(pattern)

  if len(filtered_plists) == 0 {
    fmt.Printf("No matches for '%s'\n", pattern)
    os.Exit(1)
  }

  if len(filtered_plists) > 1 {
    fmt.Println("Too many matches:")

    for k, _ := range filtered_plists {
      fmt.Println(k)
    }

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
    // Check if the dir exists. If it doesn't, move on to the next. Don't throw
    // an error.
    if _, err := os.Stat(dir); err != nil {
      continue
    }

    plist_files, err := ioutil.ReadDir(dir)
    check_error(err)

    for _, file := range plist_files {
      plists[file.Name()] = path.Join(dir, file.Name())
    }
  }

  return plists
}

// Get a list of direcotries to search.
func dirs() []string {
  dirs := user_dirs()

  if is_root() {
    root_dirs := []string{"/Library/LaunchDaemons", "/System/Library/LaunchDaemons"}
    dirs = append(dirs, root_dirs ...)
  }

  return dirs
}

func user_dirs() []string {
  user, err := user.Current()
  check_error(err)

  return []string{path.Join(user.HomeDir, "Library", "LaunchAgents"), "/Library/LaunchAgents"}
}

// Determine if the current process is running as root.
func is_root() bool {
  return os.Geteuid() == 0
}
