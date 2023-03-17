package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	workspace := prompt("Enter your Bitbucket workspace name: ")
	for workspace == "" {
		workspace = prompt("Workspace name cannot be empty. Please enter a valid workspace name: ")
	}

	username := prompt("Enter your Bitbucket username: ")
	for username == "" {
		username = prompt("Username cannot be empty. Please enter a valid username: ")
	}

	password := promptPassword()
	for password == "" {
		password = promptPassword()
	}

	action := promptAction("Do you want to (s)how the list of repositories or (c)lone them? ")
	for action != "s" && action != "c" {
		action = promptAction("Invalid option. Please choose (s)how or (c)lone: ")
	}

	perPageStr := prompt("Enter number of repositories per page (default 10): ")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		perPage = 10
	}

	pageStr := prompt("Enter page number (default 1): ")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	if action == "s" {
		showRepositories(perPage, page, workspace, username, password)
	} else {
		cloneRepositories(perPage, page, workspace, username, password)
	}
}

func prompt(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func promptPassword() string {
	fmt.Print("Enter your Bitbucket password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}
	if len(password) == 0 {
		fmt.Println("Password cannot be empty.")
		return promptPassword()
	}
	return string(password)
}

func promptAction(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	char, _, _ := reader.ReadRune()
	return strings.ToLower(string(char))
}

func getRepositories(perPage int, page int, workspace, username, password string) ([]string, error) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s?pagelen=%d&page=%d", workspace, perPage, page)

	cmd := exec.Command("curl", "-s", "-u", fmt.Sprintf("%s:%s", username, password), url)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing curl command: %s", err.Error())
	}

	type repo struct {
		Slug string `json:"slug"`
	}

	type response struct {
		Values []repo `json:"values"`
	}

	var data response
	err = json.Unmarshal(output, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response data: %s", err.Error())
	}

	var repositories []string
	for _, r := range data.Values {
		repositories = append(repositories, r.Slug)
	}

	return repositories, nil
}

func showRepositories(perPage int, page int, workspace string, username string, password string) {
	repositories, err := getRepositories(perPage, page, workspace, username, password)
	if err != nil {
		fmt.Printf("Error getting repositories: %s\n", err.Error())
		return
	}
	for _, repo := range repositories {
		fmt.Println(repo)
	}
}

func cloneRepositories(perPage int, page int, workspace string, username string, password string) {
	repositories, err := getRepositories(perPage, page, workspace, username, password)
	if err != nil {
		fmt.Printf("Error getting repositories: %s\n", err.Error())
		return
	}
	for _, repo := range repositories {
		cloneRepository(workspace, repo, username, password)
	}
}

func cloneRepository(workspace string, repository string, username string, password string) {
	url := fmt.Sprintf("https://%s:%s@bitbucket.org/%s/%s.git", username, password, workspace, repository)
	cmd := exec.Command("git", "clone", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Error cloning repository %s: %s\n", repository, exitErr.Stderr)
		} else {
			fmt.Printf("Error cloning repository %s: %v\n", repository, err)
		}
	} else {
		fmt.Printf("Repository %s cloned successfully.\n", repository)
	}
}
