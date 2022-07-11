package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ITA-Dnipro/Dp-230-Test-LFI/internal/client"
	"github.com/ITA-Dnipro/Dp-230-Test-LFI/internal/lfiscanner"
)

const (
	DirLevelUpAttempts = 10
)

var Targets = map[string]string{
	"/etc/passwd": "root:x",
}

func main() {
	cl := client.New()
	lfiScanner := lfiscanner.New(&lfiscanner.Config{
		Targets: Targets,
		LevelUpAttempts: 10,
	}, cl)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		results, err := lfiScanner.ScanUrl(sc.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to scan url: %s\n", err.Error())
			continue
		}

		fmt.Fprintln(os.Stdin, results)
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error while scanning input:", err)
	}
}
