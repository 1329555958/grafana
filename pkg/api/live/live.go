package live

import (
	"net/http"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/middleware"
)

type LiveConn struct {
	log log.Logger
}

func New() *LiveConn {
	go h.run()

	return &LiveConn{log: log.New("live.server")}
}

func (lc *LiveConn) Serve(w http.ResponseWriter, r *http.Request) {
	lc.log.Info("Upgrading to WebSocket")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(3, "Live: Failed to upgrade connection to WebSocket", err)
		return
	}

	c := newConnection(ws)
	h.register <- c

	go c.writePump()
	c.readPump()
}

func (lc *LiveConn) PushToStream(c *middleware.Context, message dtos.StreamMessage) {
	h.streamChannel <- &message
	c.JsonOK("Message recevived")
}
