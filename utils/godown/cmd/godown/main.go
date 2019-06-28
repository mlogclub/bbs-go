package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mlogclub/mlog/utils/godown"
)

var (
	guesslang = flag.String("g", "", "guesslang")
	option    *godown.Option
)

func guesslanger(code string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", *guesslang)
	} else {
		cmd = exec.Command("sh", "-c", *guesslang)
	}
	cmd.Stdin = strings.NewReader(code)
	b, err := cmd.CombinedOutput()
	return strings.ToLower(strings.TrimSpace(string(b))), err
}

func main() {
	flag.Parse()
	if *guesslang != "" {
		option = &godown.Option{GuessLang: guesslanger}
	}
	option := &godown.Option{GuessLang: guesslanger}
	if err := godown.Convert(os.Stdout, os.Stdin, option); err != nil {
		log.Fatal(err)
	}
}
