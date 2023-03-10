package server

import (
	zlog "github.com/rs/zerolog/log"
)

func Server(command string) {
	switch command := command; command {
	case "start":
		zlog.Print(command)
	case "stop":
		zlog.Print(command)
	case "save":
		zlog.Print(command)
	}
}
