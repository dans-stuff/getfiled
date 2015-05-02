package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	boot()
	wait()
}

func boot() {
	flag.Parse()
	parts := flag.Args()
	if len(parts) != 2 {
		log.Fatalf("Syntax: getfiled file generator\n")
	}
	startGenerating(parts[0], parts[1])
}
func wait() {
	sigs := make(chan os.Signal, 0)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	sig := <-sigs
	log.Println("Got signal", sig, ", quitting.")
}

func startGenerating(fname, cmdString string) {
	parts := strings.Fields(cmdString)
	go genFile(fname, parts)
}

func genFile(fname string, parts []string) {
	num := 0
	log.Printf("Creating `%s` to provide `%s`\n", parts, fname)

	createFifo := func() {
		fi, statErr := os.Stat(fname)
		if statErr == nil && (fi.Mode()|os.ModeNamedPipe == 0) {
			log.Fatalf("%s exists and is not a FIFO, quitting for safety\n", fname)
		}
		tmp := fmt.Sprintf("%s%d", fname, num)
		num += 1
		if creErr := syscall.Mknod(tmp, syscall.S_IFIFO|0666, 0); creErr != nil {
			log.Fatalf("Failed to create fifo %s: %s", tmp, creErr.Error())
		}
		os.Rename(tmp, fname)
	}

	createFifo()
	for {
		writeFile, openErr := os.OpenFile(fname, os.O_WRONLY, os.ModeNamedPipe|os.ModePerm|os.ModeExclusive)
		createFifo()
		if openErr != nil {
			log.Fatalf("Failed to feed %s: %s", fname, openErr.Error())
		}
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = writeFile

		go func() {
			if execErr := cmd.Run(); execErr != nil {
				log.Printf("Failed to run %s because %s\n", parts, execErr.Error())
			}
			writeFile.Close()
			log.Printf("Executed `%s` to provide `%s`", parts, fname)
		}()
	}
}
