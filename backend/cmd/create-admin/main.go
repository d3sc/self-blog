package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/personal-blog/cms/backend/internal/auth"
	"github.com/personal-blog/cms/backend/internal/config"
	"golang.org/x/crypto/ssh/terminal"
)

func readPassword(prompt string, reader *bufio.Reader) ([]byte, error) {
	fmt.Print(prompt)
	if terminal.IsTerminal(int(syscall.Stdin)) {
		value, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		return value, err
	}
	value, err := reader.ReadString('\n')
	fmt.Println()
	return []byte(strings.TrimRight(value, "\r\n")), err
}

func main() {
	u := flag.String("username", "", "admin username")
	flag.Parse()
	if *u == "" || len(*u) > 100 || strings.TrimSpace(*u) != *u {
		fmt.Fprintln(os.Stderr, "valid -username is required")
		os.Exit(2)
	}
	cfg, e := config.Load()
	if e != nil { fmt.Fprintln(os.Stderr, e); os.Exit(1) }
	reader := bufio.NewReader(os.Stdin)
	p, e := readPassword("Password: ", reader); if e != nil { fmt.Fprintln(os.Stderr, e); os.Exit(1) }
	q, e := readPassword("Confirm password: ", reader); if e != nil { fmt.Fprintln(os.Stderr, e); os.Exit(1) }
	if string(p) != string(q) || len(p) < 12 { fmt.Fprintln(os.Stderr, "passwords must match and be at least 12 characters"); os.Exit(2) }
	h, e := auth.HashPassword(string(p)); if e != nil { fmt.Fprintln(os.Stderr, e); os.Exit(1) }
	db, e := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if e != nil { fmt.Fprintln(os.Stderr, e); os.Exit(1) }
	defer db.Close()
	_, e = db.Exec(context.Background(), "INSERT INTO admins(username,password_hash) VALUES($1,$2)", *u, h)
	if e != nil { fmt.Fprintln(os.Stderr, e); os.Exit(1) }
	fmt.Println("admin created")
}
