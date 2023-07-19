package sockmon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Version = "dev"

var cache map[string]Socket
var dumpFilename string
var errFilename string
var bindAddress string
var dsn string
var db *gorm.DB
var filter string
var configFile string

const CACHE_SIZE int = 10000

var cmd *cobra.Command = &cobra.Command{
	Use:     "sockmon",
	RunE:    fn,
	Version: Version,
}

func NewCommand() *cobra.Command {
	return cmd
}

func fn(cmd *cobra.Command, args []string) error {
	zapLog, _ := zap.NewDevelopment()
	log := zapLog.Sugar()
	zap.ReplaceGlobals(zapLog)

	configFile, _ = cmd.PersistentFlags().GetString("config")
	viper.SetConfigFile(configFile)
	viper.ReadInConfig()

	// Retrieve config from file or use defaults from command line flags
	dumpFilename = viper.GetString("dump-file")
	errFilename = viper.GetString("error-file")
	bindAddress = viper.GetString("bind-address")
	dsn = viper.GetString("postgres")
	filter = viper.GetString("filter")

	cache = make(map[string]Socket, CACHE_SIZE)

	if dsn != "" {
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database. err: ", err)
		}
		if err := db.AutoMigrate(&SockmonStat{}); err != nil {
			log.Fatal("Failed DB initial migration. err: ", err)
		}

	}
	go func() {
		http.HandleFunc("/", handlerDefault)
		http.HandleFunc("/rtt", handlerRttOnly)
		if err := http.ListenAndServe(bindAddress, nil); err != nil {
			log.Fatalf("Invalid bind address. err: %s", err)
		}
	}()
	cmdName := "stdbuf"
	cmdArgs := []string{"-i0", "-o0", "-e0", "ss", "-ntieEOH"}
	if filter != "" {
		cmdArgs = append(cmdArgs, filter)
	}

	ec := exec.Command(cmdName, cmdArgs...)
	stdout, err := ec.StdoutPipe()
	if err != nil {
		return err
	}
	if err := ec.Start(); err != nil {
		log.Errorf("%v\n", err)
		return err
	}
	log.Infof("sockmon starting.")
	for k, v := range viper.AllSettings() {
		log.Debugf("\t %s -> %s", k, v)
	}

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
		log.Errorf("Invalid ss input. err: %s", err)
	} else {
		if !isValidOutput(in, sock) {
			log.Errorf("Invalid output. -> %s", sock.Key())
		}
		// to local memory cache
		cacheStore(sock)
		// to log file
		if dumpFilename != "" {
			// By default, it does not dump to file.
			if err := ssDumpFile(sock); err != nil {
				log.Fatalf("Cannot dump to file. err: %s", err)
			}
		}
		// to DB
		if dsn != "" {
			log.Infof("create: %s", sock)
			stat := SockmonStat{
				Timestamp:                 sock.Timestamp,
				Src:                       sock.Src,
				Dst:                       sock.Dst,
				Protocol:                  sock.Protocol,
				Sport:                     sock.Sport,
				Dport:                     sock.Dport,
				SocketExtendedInformation: sock.Ext,
			}
			statRes := db.Create(&stat)
			if statRes.Error != nil {
				log.Errorf("DB update error, socket: %s", sock)
			}
		}
	}
}

func handlerDefault(w http.ResponseWriter, r *http.Request) {
	// filter by query parameters
	params := r.URL.Query()
	c := FilterByParams(params, cache)
	out, err := json.MarshalIndent(&c, "", "  ")

	if err != nil {
		io.WriteString(w, fmt.Sprintf("{'err':'%s'}\n", err.Error()))
		return
	}
	io.WriteString(w, fmt.Sprintf("%s\n", string(out)))
}

func handlerRttOnly(w http.ResponseWriter, r *http.Request) {
	// filter by query parameters
	params := r.URL.Query()
	c := FilterByParams(params, cache)

	rttcache := map[string]float32{}
	for key, val := range c {
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
	cmd.PersistentFlags().StringP("config", "c", "", "Use: sockmon --config <CONFIG_PATH> or sockmon -c <CONFIG_PATH>  Various file formats such as YML, TOML and JSON are available.")

	cmd.PersistentFlags().StringP("dump-file", "d", "", "Use: sockmon --dump-file <FILENAME> or sockmon -d <FILENAME> (by default, it does not dump to file.) ")
	cmd.PersistentFlags().StringP("error-file", "e", "", "Use: sockmon --error-file <FILENAME> or sockmon -e <FILENAME> (by default, it does not dump to file.) ")
	cmd.PersistentFlags().StringP("bind-address", "b", ":8931", "Use: sockmon --bind-address <Address:Port> or sockmon -b <Address:Port> ")
	cmd.PersistentFlags().StringP("postgres", "p", "", "Use: sockmon --postgres 'postgres://user:password@localhost:5432/dbname' or sockmon -p 'postgres://user:password@localhost:5432/dbname' ")
	cmd.PersistentFlags().StringP("filter", "f", "", "Use: sockmon --filter '<FILTER>' or sockmon -f '<FILTER>' ss filter.  Please take a look at the iproute2 official documentation. e.g. dport = :80 ")

	viper.BindPFlag("dump-file", cmd.PersistentFlags().Lookup("dump-file"))
	viper.BindPFlag("error-file", cmd.PersistentFlags().Lookup("error-file"))
	viper.BindPFlag("bind-address", cmd.PersistentFlags().Lookup("bind-address"))
	viper.BindPFlag("postgres", cmd.PersistentFlags().Lookup("postgres"))
	viper.BindPFlag("filter", cmd.PersistentFlags().Lookup("filter"))
}
