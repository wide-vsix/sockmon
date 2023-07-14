package sockmon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cache map[string]Socket
var dumpFilename string
var errFilename string
var bindAddress string

const CACHE_SIZE int = 10000

var cmd *cobra.Command = &cobra.Command{
	Use:  "sockmon",
	RunE: fn,
}

func NewCommand() *cobra.Command {
	return cmd
}

func fn(cmd *cobra.Command, args []string) error {
	zapLog, _ := zap.NewDevelopment()
	log := zapLog.Sugar()
	zap.ReplaceGlobals(zapLog)

	dumpFilename, _ = cmd.PersistentFlags().GetString("dump-file")
	errFilename, _ = cmd.PersistentFlags().GetString("error-file")
	bindAddress, _ = cmd.PersistentFlags().GetString("bind-address")
	cache = make(map[string]Socket, CACHE_SIZE)

	go func() {
		http.HandleFunc("/", handlerDefault)
		http.HandleFunc("/rtt", handlerShort)
		if err := http.ListenAndServe(bindAddress, nil); err != nil {
			log.Fatalf("Invalid bind address. err: %s", err)
		}
	}()

	ec := exec.Command("stdbuf", "-i0", "-o0", "-e0", "ss", "-ntieEOH")
	stdout, err := ec.StdoutPipe()
	if err != nil {
		return err
	}
	if err := ec.Start(); err != nil {
		log.Errorf("%v\n", err)
		return err
	}

	log.Infof("starting. dump-file: '%s' bind-address'%s'", dumpFilename, bindAddress)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		input(scanner.Text())
	}

	if err := ec.Wait(); err != nil {
		log.Error("Error at ss process.", err)
		return err
	}
	return nil
}

func input(in string) {
	sock, err := ParseSsOutput(in)
	if err != nil {
		log.Fatalf("Invalid ss input. err: %s", err)
	}
	cacheStore(sock)
	if dumpFilename != "" {
		// By default, it does not dump to file.
		if err := ssDumpFile(sock); err != nil {
			log.Fatalf("Cannot dump to file. err: %s", err)
		}
	}
	if !isValidOutput(in, sock) {
		log.Errorf("Invalid output. -> %s", sock.Key())
	}
}

func handlerDefault(w http.ResponseWriter, r *http.Request) {
	out, err := json.MarshalIndent(&cache, "", "  ")
	if err != nil {
		io.WriteString(w, fmt.Sprintf("{'err':'%s'}\n", err.Error()))
		return
	}
	io.WriteString(w, fmt.Sprintf("%s\n", string(out)))
}

func handlerShort(w http.ResponseWriter, r *http.Request) {
	rttcache := map[string]float32{}
	for key, val := range cache {
		rttcache[key] = val.Ext.Rtt
	}
	out, err := json.MarshalIndent(&rttcache, "", "  ")
	if err != nil {
		io.WriteString(w, fmt.Sprintf("{'err':'%s'}\n", err.Error()))
		return
	}
	io.WriteString(w, fmt.Sprintf("%s\n", string(out)))
}

func init() {
	cmd.PersistentFlags().String("dump-file", "", "Use: sockmon --dump-file <FILENAME> (by default, it does not dump to file.) ")
	cmd.PersistentFlags().String("error-file", "", "Use: sockmon --error-file <FILENAME> (by default, it does not dump to file.) ")
	cmd.PersistentFlags().String("bind-address", ":8931", "Use: sockmon --bind-address <Address:Port> ")
}