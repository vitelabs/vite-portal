package rpc

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

func StartWsRpc(port int32, timeout int64) {
	nodeHub := ws.NewHub(timeout, messageHandler)
	go nodeHub.Run()

	relayerHub := ws.NewHub(timeout, messageHandler)
	go relayerHub.Run()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ws/v1/node", func(w http.ResponseWriter, r *http.Request) {
		handleNode(nodeHub, w, r, timeout)
	})
	serveMux.HandleFunc("/ws/v1/relayer", func(w http.ResponseWriter, r *http.Request) {
		handleRelayer(relayerHub, w, r, timeout)
	})
	serveMux.HandleFunc("/", home)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), serveMux)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("WS RPC closed")
		os.Exit(1)
	}
}

func messageHandler(client *ws.Client, msg []byte) {
	logger.Logger().Info().Msg(fmt.Sprintf("%s", msg))
	client.Send <- msg
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws/v1/node?chain=vite_testnet")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
