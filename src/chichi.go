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
	if args.ContainsArg(os.Args, "-h", "--help") {
		printHelp()
		os.Exit(0)
	}
	if args.ContainsArg(os.Args, "-v", "--version") {
		fmt.Printf("chichi %s\nhttps://github.com/tvrzna/chichi\n\nReleased under the MIT License.\n", version)
		os.Exit(0)
	}

	conf := loadConfig(filepath.Join(getConfigPath(), configFileName), os.Args)

	go loop(Normal, conf.ShortPeriod, conf.ShortBreak, "short")
	go loop(Critical, conf.LongPeriod, conf.LongBreak, "looong")

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

func loop(level UrgencyLevel, periodLength int, breakLength int, breakType string) {
	for {
		time.Sleep(time.Duration(periodLength) * time.Second)
		message := formatBreakMessage(breakLength, breakType)
		n := &SendNotify{urgencyLevel: level, length: breakLength, message: message}
		log.Printf("time for %s break for %d seconds", breakType, breakLength)
		if err := n.Send(); err != nil {
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
	levels := []string{"seconds", "minutes", "hours", "days"}
	// levelsTypes := []int{60, 60, 24}

	level := 0
	length := breakLength
	// for i := 0; i < len(levelsTypes); i++ {
	// 	if (length % levelsTypes[i]) >= 1 {
	// 		length = length / levelsTypes[i]
	// 		level = i
	// 	}
	// }

	breakEnd := time.Now().Add(time.Duration(breakLength) * time.Second)

	return fmt.Sprintf("It's a time for %s break for %d %s.\n\nThe break ends at %s.", breakType, length, levels[level], breakEnd.Format("15:04:05"))
}

func printHelp() {
	fmt.Println("TODO!")
}

func getVersion() string {
	if buildVersion != "" {
		return buildVersion[1:]
	}
	return version
}
