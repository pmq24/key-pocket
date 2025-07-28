package log

import "log"

var Verbose bool

func Verboseln(msg string) {
	if (Verbose) {
		log.Println(msg)
	}
}

func Verbosef(msg string, a ...any) {
	if (Verbose) {
		log.Printf(msg, a...)
	}
}

func Infoln(msg string) {
  log.Println(msg)
}

func Infof(msg string, a ...any) {
	log.Printf(msg, a...)
}

func Warnln(msg string) {
	log.Println("WARN: " + msg)
}

func Warnf(msg string, a ...any) {
	log.Printf("WARN: " + msg, a...)
}

func Errorln(msg string) {
	log.Println("ERROR: " + msg)
}

func Errorf(msg string, a ...any) {
	log.Printf("ERROR: " + msg, a...)
}
