package server

import "fmt"

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiGreen  = "\033[32m"
	ansiBlue   = "\033[34m"
	ansiCyan   = "\033[36m"
	ansiYellow = "\033[33m"
)

func (s *Server) printBootstrapStatus() {
	fmt.Printf(
		"%s%s[boot]%s %sDatabase connected%s\n",
		ansiBold,
		ansiGreen,
		ansiReset,
		ansiCyan,
		ansiReset,
	)
}

func (s *Server) printStartupBanner() {
	address := "http://localhost:" + s.Config.AppPort

	fmt.Printf(
		"\n%s+--------------------------------------------------+%s\n",
		ansiBlue,
		ansiReset,
	)
	fmt.Printf(
		"%s|%s %sAUTO FLOW API%s                               %s|%s\n",
		ansiBlue,
		ansiReset,
		ansiBold,
		ansiReset,
		ansiBlue,
		ansiReset,
	)
	fmt.Printf(
		"%s|%s %sDB:%s     connected                           %s|%s\n",
		ansiBlue,
		ansiReset,
		ansiGreen,
		ansiReset,
		ansiBlue,
		ansiReset,
	)
	fmt.Printf(
		"%s|%s %sHTTP:%s   %s%-35s%s %s|%s\n",
		ansiBlue,
		ansiReset,
		ansiGreen,
		ansiReset,
		ansiYellow,
		address,
		ansiReset,
		ansiBlue,
		ansiReset,
	)
	fmt.Printf(
		"%s+--------------------------------------------------+%s\n\n",
		ansiBlue,
		ansiReset,
	)
}
