# Wake-On-Lan Server

This is a minimalistic server to wake up computers in your network using a web-gui

- Targets are identified by their IP and MAC.
- The server pretends to run PHP. 
- Network status identification works via ping and connection attempts to different ports.

See [Docker/README.md] if you want to build & run in Docker

## build
with go already installed

    GOPATH=$PWD go build src/wakeUp.go


## run
create a `config.json` according to your environment and just do it (✔)

    ./wakeUp

if you want to serve a favicon place a file "favicon.ico" in the template dir

## usage
`-port string`	port number (default "8000")
`-root` 		run in root mode (default true)
`-v`			verbose

## endpoints
the server supports the following endpoints:
- `/favicon.ico` responds with the content of the template/favicon.ico if present
- `/wake_up.php` uses GET['id'] to wakeup the coresponding entry
- `/netstat.php` uses GET['id'] to respond with the determined state (offline/online)
- `/index.php` (and every other endpoint) shows `status.html`

## template
status.html is a template for the server, to change it you might want to know:

    {{range $id, $ele := .}}
    …
    {{end}}

iterates over all entries in `config.json`

`{{$ele.Name}}` displays the `Name` entry
`{{$ele.Text}}` displays the `Text` entry
`{{$id}}` displays the internal id of the entry (the position starting with zero)

further information about the templating engine can be found here: [https://golang.org/pkg/html/template]

## compression
I like to compress the resulting template just for fun:

	(
		cat template/status.html | tr -d '\n' | sed 's/<script>.*<\/script>.*/<script>/' # before script tag
		cat template/status.html | tr -d '\n' | sed 's/.*<script>//;s/<\/script>.*//' | uglifyjs --mangle | sed 's/;$//' # compress script
		cat template/status.html | tr -d '\n' | sed 's/.*<script>.*<\/script>/<\/script>/' # after script tag
	) | html-minifier \
		--collapse-boolean-attributes \
		--collapse-inline-tag-whitespace \
		--collapse-whitespace \
		--decode-entities \
		--minify-css \
		--minify-js \
		--prevent-attributes-escaping \
		--remove-attribute-quotes \
		--remove-comments \
		--remove-empty-attributes \
		--remove-optional-tags \
		--remove-redundant-attributes \
		--remove-script-type-attributes \
		--remove-style-link-type-attributes \
		--remove-tag-whitespace \
		--use-short-doctype \
	> template/status.min.html