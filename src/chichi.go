package src

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"github.com/tvrzna/go-utils/args"
)

const version = "0.1.0"

var buildVersion string

const (
	xdgConfigHome  = "XDG_CONFIG_HOME"
	configFileName = "chichi.conf"
)

func Run() {
	handleServiceArgs(os.Args)

	conf, err := loadConfig(filepath.Join(getConfigPath(), configFileName), os.Args)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go loop(conf, Normal, conf.ShortPeriod, conf.ShortBreak, "short")
	go loop(conf, Critical, conf.LongPeriod, conf.LongBreak, "looong")

	wait()
}

func getConfigPath() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		configPath := os.Getenv(xdgConfigHome)
		if configPath == "" {
			configPath = filepath.Join(configPath, configFileName)
		} else {
			user, _ := user.Current()
			configPath = filepath.Join(user.HomeDir, ".config", configFileName)
		}

	}
	return configPath
}

func loop(conf *Config, level UrgencyLevel, periodLength int, breakLength int, breakType string) {
	for {
		time.Sleep(time.Duration(periodLength) * time.Second)
		message := formatBreakMessage(breakLength, breakType)
		n := &SendNotify{urgencyLevel: level, length: breakLength, message: message}
		log.Printf("time for %s break for %d seconds", breakType, breakLength)
		if err := n.Send(conf.NotifySendPath); err != nil {
			log.Println(err)
		}
	}
}

func wait() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM)
	<-c
}

func formatBreakMessage(breakLength int, breakType string) string {
	breakEnd := time.Now().Add(time.Duration(breakLength) * time.Second)
	length, unit := formatBreakLength(breakLength)

	return fmt.Sprintf("It's a time for %s break for %d %s.\n\nThe break ends at %s.", breakType, length, unit, breakEnd.Format("15:04:05"))
}

func formatBreakLength(breakLength int) (length int, unit string) {
	levels := []string{"second", "minute", "hour", "day"}
	levelsTypes := []int{60, 60, 24}

	level := 0
	length = breakLength
	for i := 0; i < len(levelsTypes); i++ {
		if (length / levelsTypes[i]) >= 1 {
			length = length / levelsTypes[i]
			level = i + 1
		} else {
			break
		}
	}
	unit = levels[level]
	if length > 1 {
		unit += "s"
	}

	return length, unit
}

func handleServiceArgs(arguments []string) {
	if args.ContainsArg(arguments, "-h", "--help") {
		printHelp()
		os.Exit(0)
	}
	if args.ContainsArg(arguments, "-v", "--version") {
		fmt.Printf("chichi %s\nhttps://github.com/tvrzna/chichi\n\nReleased under the MIT License.\n", getVersion())
		os.Exit(0)
	}
}

func printHelp() {
	fmt.Println("Usage: chichi [options]")
	fmt.Println("Options:")
	fmt.Printf("  -h, --help\t\t\tprint this help\n")
	fmt.Printf("  -v, --version\t\t\tprint version\n")
	fmt.Printf("  -sp, --short-period TIME\tperiod length between short breaks\n")
	fmt.Printf("  -sb, --short-break TIME\tlength of short breaks\n")
	fmt.Printf("  -lp, --long-period TIME\tperiod length between long breaks\n")
	fmt.Printf("  -lb, --long-break TIME\tlength of long breaks\n")
}

func getVersion() string {
	if buildVersion != "" {
		return buildVersion[1:]
	}
	return version
}
