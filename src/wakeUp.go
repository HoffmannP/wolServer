package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"networking"
	"os"
	"sort"
	"strconv"
)

type nodes []struct {
	Name string
	Text string
	IP   string
	MAC  string
}

type mainObject struct {
	r bool
	l bool
	p string
	t *template.Template
	n nodes
}

func (n *nodes) Len() int {
	return len(*n)
}
func (n *nodes) Less(i, j int) bool {
	return (*n)[i].Name < (*n)[j].Name
}
func (n *nodes) Swap(i, j int) {
	(*n)[i], (*n)[j] = (*n)[j], (*n)[i]
	return
}

func status(ip string) string {
	return "offline"
}

func main() {
	me := new(mainObject)

	// logging
	flag.BoolVar(&me.r, "root", true, "run in root mode")
	flag.BoolVar(&me.l, "v", false, "verbose")
	flag.StringVar(&me.p, "port", "8000", "port number")
	flag.Parse()
	if me.l {
		fmt.Printf("Rootmode: %t\n", me.r)
	}

	// read configuration
	config, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(config)
	err = dec.Decode(&me.n)
	sort.Sort(&me.n)
	if err != nil {
		panic(err)
	}

	// read template
	me.t, err = template.New("status.html").Funcs(template.FuncMap{"status": status}).ParseFiles("template/status.html")
	if err != nil {
		panic(err)
	}

	// start server
	http.Handle("/", me)
	if me.l {
		fmt.Printf("Running WOL-server on port %s\n", me.p)
	}
	http.ListenAndServe(fmt.Sprintf(":"+me.p), nil)
}

func (me *mainObject) fav(w http.ResponseWriter) {
	f, err := os.Open("template/favicon.ico")
	if err != nil {
		fmt.Fprint(w, "")
		return
	}
	_, err = io.Copy(w, f)
	if err != nil {
		panic(err)
	}
}

func (me *mainObject) wol(ID int) {
	Broadcast := "255.255.255.255"
	err := networking.SendMagicPacket(me.n[ID].MAC, Broadcast)
	if err != nil {
		fmt.Println(err)
	}
	// re-send package to IP of target
	err = networking.SendMagicPacket(me.n[ID].MAC, me.n[ID].IP)
	if err != nil {
		fmt.Println(err)
	}
	if me.l {
		println("Sending magic packet to", me.n[ID].Name)
	}
}

func (me *mainObject) netstat(w http.ResponseWriter, ID int) {
	var status bool
	var err error
	if me.r {
		status, err = networking.Ping(me.n[ID].IP)
	} else {
		status, err = networking.Netstats(me.n[ID].IP, []string{"135", "137", "138", "139", "445", "593", "3389", "22"})
	}
	if err != nil {
		fmt.Println(err)
		status = false
	}
	if me.l {
		fmt.Printf("  %-25s online: %t\n", me.n[ID].Name, status)
	}
	if status {
		fmt.Fprint(w, "<span class='online'>online</span>")
	} else {
		fmt.Fprint(w, "<span class='offline'>offline</span>")
	}
}

func (me *mainObject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		me.fav(w)
		return
	}

	if r.URL.Path == "/wake_up.php" {
		r.ParseForm()
		id, err := strconv.Atoi(r.Form.Get("ID"))
		if err == nil {
			me.wol(id)
		}
		http.Redirect(w, r, "/index.php", 302)
		return
	}

	if r.URL.Path == "/netstat.php" {
		r.ParseForm()
		id, err := strconv.Atoi(r.Form.Get("ID"))
		if err == nil {
			me.netstat(w, id)
		}
		return
	}

	if r.URL.Path != "/index.php" {
		http.Redirect(w, r, "/index.php", 302)
		return
	}

	err := me.t.Execute(w, me.n)
	if err != nil {
		panic(err)
	}
}
