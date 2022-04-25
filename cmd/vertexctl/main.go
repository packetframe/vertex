package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/packetframe/vertex/internal/db"
)

// Set by build process
var version = "dev"

// Flags
var verbose bool

// Environment variables
var server string

func init() {
	cobra.OnInitialize(func() {
		server = os.Getenv("VERTEX_SERVER")
		if server == "" {
			log.Warn("VERTEX_SERVER environment variable not set, setting to http://localhost:8080")
			server = "http://localhost:8080"
		}
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
	})
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.MarkFlagRequired("server")
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(findCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(versionCmd)
}

var rootCmd = &cobra.Command{
	Use: "vertexctl",
}

var createCmd = &cobra.Command{
	Use:   "create [filter] [expire] [name]",
	Short: "create a new license",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			log.Fatal("Usage: vertexctl create [filter] [expire] [name]")
		}

		resp := req("/rules", http.MethodPost, map[string]string{"filter": args[0], "expire": args[1], "name": args[2]})
		defer resp.Body.Close()

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Fatal(err)
		}
		log.Info(response["message"])
	},
}

var findCmd = &cobra.Command{
	Use:     "list",
	Short:   "list rules",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		resp := req("/rules", http.MethodGet, nil)
		defer resp.Body.Close()

		var b struct {
			Data struct {
				Rules []db.Rule `json:"rules"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
			log.Fatal(err)
		}

		var rows []table.Row
		for _, rule := range b.Data.Rules {
			rows = append(rows, table.Row{rule.ID, rule.Name, rule.ExpireStr, rule.CreatedAt, rule.Filter})
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{fmt.Sprintf("ID (total %d)", len(rows)), "Name", "Expire", "Created", "Filter"})
		t.AppendRows(rows)
		t.AppendSeparator()
		t.SetStyle(table.StyleLight)
		t.Render()
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete a rule",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Usage: vertexctl delete [id]")
		}

		resp := req("/rules/"+args[0], http.MethodDelete, nil)
		defer resp.Body.Close()

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Fatal(err)
		}
		log.Info(response["message"])
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version information",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("vertexctl version %s server %s", version, server)
	},
}

func req(path, method string, params map[string]string) *http.Response {
	u, err := url.Parse(server)
	if err != nil {
		log.Fatal(err)
	}
	u.Path = path
	q := u.Query()
	u.RawQuery = q.Encode()
	log.Debugf("Connecting to %s", u.String())
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		log.Fatal(string(respBody))
	}
	return resp
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
