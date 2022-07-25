package stream

import (
	"time"

	"net/http"

	"github.com/gorilla/websocket"
	logger "github.com/kispeagle/go-binance/logs"
)

func Subscribe(url string, c chan interface{}) error {

	header := make(http.Header)
	header.Add("method", "SUBSCRIBE")
	header.Add("params", "btcusdt@aggTrade")
	header.Add("params", "btcusdt@depth")
	header.Add("id", "1")

	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		logger.Log.Error(err)
		c <- nil
		return err
	}
	defer conn.Close()

	ticker := time.NewTicker(PingPeriod)

	for {
		select {
		// ping check connection
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(60))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
		default:
			// read data from socket
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
			logger.Log.Debug(string(msg))

		}
	}

	return nil
}
