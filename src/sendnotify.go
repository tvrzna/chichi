package src

import (
	"os/exec"
	"strconv"
)

// Urgency Level for send-notify command
type UrgencyLevel int

const (
	Low UrgencyLevel = iota
	Normal
	Critical
)

func (u UrgencyLevel) String() string {
	return []string{"low", "normal", "critical"}[u]
}

type SendNotify struct {
	urgencyLevel UrgencyLevel
	length       int
	message      string
}

func (s *SendNotify) strLength() string {
	return strconv.Itoa(1000 * s.length)
}

func (s *SendNotify) Send(path string) error {
	_, err := exec.Command(path, "-u", s.urgencyLevel.String(), "-t", s.strLength(), s.message).Output()
	return err
}
