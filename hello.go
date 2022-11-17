package main

import (
	"C"
	"bufio"
	"log"

	bpf "github.com/aquasecurity/libbpfgo"
)
import (
	"fmt"
	"os"
	"os/signal"
)

func TracePrint() {
	f, err := os.Open("/sys/kernel/debug/tracing/trace_pipe")
	if err != nil {
		fmt.Println("TracePrint failed to open trace pipe: %v", err)
		return
	}
	r := bufio.NewReader(f)
	b := make([]byte, 1000)
	for {
		len, err := r.Read(b)
		if err != nil {
			fmt.Println("TracePrint failed to read from trace pipe: %v", err)
			return
		}
		s := string(b[:len])
		fmt.Println(s)
	}
}

func checkEnvPath(env string) (string, error) {
	filePath, _ := os.LookupEnv(env)
	if filePath == "" {
		return "", nil
	}
	_, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open %s: %w", filePath, err)
	}
	return filePath, nil
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	bpfObjPath := "hello.bpf.o"
	newModuleArgs := bpf.NewModuleArgs{
		BPFObjPath: bpfObjPath,
	}

	btfFilePath, err := checkEnvPath("BTF_FILE")
	if err != nil {
		log.Fatal(err)
	}
	if btfFilePath != "" {
		newModuleArgs.BTFObjPath = btfFilePath
	}

	bpfModule, err := bpf.NewModuleFromFileArgs(newModuleArgs)
	must(err)

	defer bpfModule.Close()

	err = bpfModule.BPFLoadObject()
	must(err)

	iter := bpfModule.Iterator()
	for {
		prog := iter.NextProgram()
		if prog == nil {
			break
		}
		_, err := prog.AttachGeneric()
		must(err)
	}

	go TracePrint()

	e := make(chan []byte, 300)
	p, err := bpfModule.InitPerfBuf("events", e, nil, 1024)
	must(err)

	p.Start()

	counter := make(map[string]int, 350)
	go func() {
		for data := range e {
			comm := string(data)
			counter[comm]++
		}
	}()

	<-sig
	p.Stop()
	for comm, n := range counter {
		fmt.Printf("%s: %d\n", comm, n)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
