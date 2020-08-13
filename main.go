package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", users)
	http.ListenAndServe(":9990", nil)
}

func users(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("status.log")
	if err != nil {
		fmt.Fprintf(w, "Failed to read users: %v", err)
	}
	defer file.Close()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		if i <= 5 {
			continue
		}
		if len(scanner.Text()) == 12 {
			fmt.Fprintln(w, `</table>`)
			break
		}
		if i == 6 {
			fmt.Fprintln(w, printHeaders())
		}
		fmt.Fprintln(w, printUser(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(w, "Failed to scan users: %v", err)
	}
}

const sep string = `	 	|		 `

func printHeaders() string {
	return `
<table>
	<tr>
		<td>
			USUARIO
		</td>
		<td>
			:3010
		</td>

		<td>
			:3000
		</td>

		<td>
			:8080
		</td>

		<td>
			ÚLTIMA VEZ
		</td>



	</tr>
`
}

func printUser(str string) string {
	splitted := strings.Split(str, ",")
	return fmt.Sprintf(`
	<tr>

		<td>
			<a href="%s">%s</a>
		</td>

		<td>
			<a href="%s">%s</a>
		</td>

		<td>
			<a href="%s">%s</a>
		</td>

		<td>
			<a href="%s">%s</a>
		</td>

		<td>
			%s
		</td>


		</tr>
`,
		"http://"+splitted[0], splitted[1],
		"http://"+splitted[0]+":3010", splitted[1],
		"http://"+splitted[0]+":3000", splitted[1],
		"http://"+splitted[0]+":8080", splitted[1],
		splitted[len(splitted)-1],
	)
}
