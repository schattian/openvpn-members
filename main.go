package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Ullaakut/nmap"
)

var frequentPortsWithAlias = map[string]string{
	"80":   "httpserve",
	"3000": "rails",
	"3010": "rails alt",
	"8080": "node",
	"5432": "postgres",
}

func frequentPorts() (ports []string) {
	for v := range frequentPortsWithAlias {
		ports = append(ports, v)
	}
	return
}

func main() {
	http.HandleFunc("/", users)
	http.ListenAndServe(":9990", nil)
}

func nmapFrequentPorts(ips ...string) map[string][]string {
	openPortsByIp := make(map[string][]string)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(ips...),
		nmap.WithPorts(frequentPorts()...),
		nmap.WithContext(ctx),
	)

	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	result, _, err := scanner.Run()
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	// Use the results to print an example output
	for i, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		for _, port := range host.Ports {
			if port.Status() == nmap.Open || port.Status() == nmap.Filtered {
				openPortsByIp[ips[i]] = append(openPortsByIp[ips[i]], strconv.Itoa(int(port.ID)))
			}
		}
	}
	return openPortsByIp
}

func users(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("status.log")
	if err != nil {
		fmt.Fprintf(w, "Failed to read users: %v", err)
	}
	defer file.Close()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintln(w, printHeaders())

	currentIP, err := myPublicIP()
	if err != nil {
		fmt.Fprintf(w, "Failed to get local public ip: %v", err)
	}
	fmt.Fprintln(w, printUser(currentIP+",server,,,"))

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		if i <= 5 {
			continue
		}
		if len(scanner.Text()) == 12 {
			break
		}
		fmt.Fprintln(w, printUser(scanner.Text()))
		fmt.Fprintln(w, "</br>")
	}
	fmt.Fprintln(w, `</table>`)
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(w, "Failed to scan users: %v", err)
	}
}

func printHeaders() string {
	return `
<table>
	<tr>
		<td>
			CLIENT
		</td>
	</tr>
`
}
func portColumn(svc, port, ip string) string {
	return fmt.Sprintf(`
		<td>
			<a href="%s">%s</a>
		</td>
	`, "http://"+ip+":"+port, svc)
}

func myPublicIP() (string, error) {
	url := "https://api.ipify.org?format=text"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

func printUser(str string) string {
	splitted := strings.Split(str, ",")
	ip, alias := splitted[0], splitted[1]
	base := fmt.Sprintf(`<tr>
		<td>
			<a href="%s">%s</a>
		</td>

	`, "http://"+ip, alias)
	openPortsByIp := nmapFrequentPorts(ip)
	for _, ports := range openPortsByIp {
		for _, port := range ports {
			portAlias := frequentPortsWithAlias[port]
			base += portColumn(portAlias, port, ip)
		}
	}
	base += "</tr>"
	return base
}
