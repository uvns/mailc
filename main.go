package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println("config error:", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cli := ParseCLI()

	// ❗没有任何参数 → 显示 help
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	// ❗参数不完整 → 显示 help
	if cli.Subject == "" {
		flag.Usage()
		os.Exit(1)
	}

	if cli.Recipient == "" && cli.ListFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := Send(cfg, cli); err != nil {
		fmt.Println("send failed:", err)
		os.Exit(1)
	}

	fmt.Println("mail sent successfully")
}
