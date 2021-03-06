// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list discovered nodes",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		// fmt.Println("list called")
		w := tabwriter.NewWriter(os.Stdout, 10, 2, 2, ' ', 0)

		db, err := sql.Open("sqlite3", "/tmp/languard.db")
		if err != nil {
			log.Fatal(err)
		}
		rows, err := db.Query("select count(*) from node where state == \"on\"")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		total_on := 0
		for rows.Next() {
			rows.Scan(&total_on)
		}
		fmt.Fprintln(w, "Total up node: ",total_on)

		rows, err = db.Query("select if_name, vlan_id, ipv4, mac, hostname, groupname, state from node")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		fmt.Fprintln(w, "DEVICE\tVLAN\tIPV4\tMAC\tHOST\tGROUP\tSTATE\tBLOCK")
		for rows.Next() {
			var if_name string
			var vlan_id int
			var ipv4 string
			var mac string
			var hostname string
			var groupname string
			var state string
			var block int
			err = rows.Scan(&if_name, &vlan_id, &ipv4, &mac, &hostname, &groupname, &state, &block)
			if err != nil {
				log.Fatal(err)
			}
			is_block := false
			if block != 0 {
				is_block = true
			}
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\t%t\n",
				if_name, vlan_id, ipv4, mac, hostname, groupname, state, is_block)
		}
		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
