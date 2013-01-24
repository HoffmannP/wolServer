package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"wol"
)

type mainObject struct {
	t *template.Template
	s []struct {
		Name string
		IP   string
		MAC  string
	}
}

func status(ip string) string {
	return "offline"
}

func main() {
	me := new(mainObject)

	// read configuration
	config, err := os.Open("config.json")
	if err != nil {
		fmt.Println("error:", err)
	}
	dec := json.NewDecoder(config)
	err = dec.Decode(&me.s)
	if err != nil {
		fmt.Println("error:", err)
	}

	// read template
	me.t, err = template.New("status.html").Funcs(template.FuncMap{"status": status}).ParseFiles("status.html")
	if err != nil {
		fmt.Println("error:", err)
	}

	// start server
	http.Handle("/", me)
	http.ListenAndServe(":8080", nil)
}

func (me *mainObject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/wake_up.php" {
		r.ParseForm()
		IP := r.Form.Get("IP")
		MAC := r.Form.Get("MAC")
		for _, s := range me.s {
			if (s.MAC == MAC) && (s.IP == IP) {
				wol.SendMagicPacket(s.MAC, s.IP)
				println("Waking Up ( IP=", s.IP, ")")
			}
		}
		http.Redirect(w, r, "/index.php", 302)
		return
	}
	if r.URL.Path != "/index.php" {
		http.Redirect(w, r, "/index.php", 302)
		return
	}
	err := me.t.Execute(w, me.s)
	if err != nil {
		fmt.Println("error:", err)
	}
}
