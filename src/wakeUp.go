package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"networking"
	"os"
)

type Config struct {
	Port  int
	Nodes []struct {
		Name string
		IP   string
		MAC  string
	}
}

type node struct {
	IP  string
	MAC string
}

type mainObject struct {
	t *template.Template
	s map[string]node
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
	var cfg Config
	err = dec.Decode(&cfg)
	if err != nil {
		fmt.Println("error:", err)
	}
	me.s = make(map[string]node, len(cfg.Nodes))
	for _, n := range cfg.Nodes {
		me.s[n.Name] = node{n.IP, n.MAC}
	}

	// read template
	me.t, err = template.New("status.html").Funcs(template.FuncMap{"status": status}).ParseFiles("template/status.html")
	if err != nil {
		fmt.Println("error:", err)
	}

	// start server
	http.Handle("/", me)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}

func (me *mainObject) wol(Name string) {
	networking.SendMagicPacket(me.s[Name].MAC, me.s[Name].IP)
	println("Waking Up ( IP=", me.s[Name].IP, ")")
}

func (me *mainObject) netstat(w http.ResponseWriter, Name string) {
	status := networking.Netstat(me.s[Name].IP)
	println("IP", me.s[Name].IP, "online", status)
	if status {
		fmt.Fprint(w, "<span class='online'>online</span>")
	} else {
		fmt.Fprint(w, "<span class='offline'>offline</span>")
	}
}

func (me *mainObject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/wake_up.php" {
		r.ParseForm()
		me.wol(r.Form.Get("Name"))
		http.Redirect(w, r, "/index.php", 302)
		return
	}

	if r.URL.Path == "/netstat.php" {
		r.ParseForm()
		me.netstat(w, r.Form.Get("Name"))
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
