package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
	"github.com/manifoldco/promptui"
)

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func PasswordPrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stderr, label+" ")
		b, _ := term.ReadPassword(int(syscall.Stdin))
		s = string(b)
		if s != "" {
			break
		}
	}
	fmt.Println()
	return s
}

func usernamePrompt() string {
	username := StringPrompt("Username:")
	return username
}

func passwordPrompt() string {
	password := PasswordPrompt("Password:")
	return password
}

func SelectPrompt(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, err
}

func YesNoPrompt(label string, def bool) bool {
    choices := "Y/n"
    if !def {
        choices = "y/N"
    }

    r := bufio.NewReader(os.Stdin)
    var s string

    for {
        fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
        s, _ = r.ReadString('\n')
        s = strings.TrimSpace(s)
        if s == "" {
            return def
        }
        s = strings.ToLower(s)
        if s == "y" || s == "yes" {
            return true
        }
        if s == "n" || s == "no" {
            return false
        }
    }
}
