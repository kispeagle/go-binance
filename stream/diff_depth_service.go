package stream

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"

	logger "github.com/kispeagle/go-binance/logs"
)

type DiffDepthService struct {
	Symbol string
}

func NewDiffDepthService(symbol string) DiffDepthService {
	s := DiffDepthService{symbol}
	return s
}

func (s DiffDepthService) Call(url string, c chan interface{}) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Log.Error(err)
		c <- nil
		return err
	}
	defer conn.Close()

	ticker := time.NewTicker(PingPeriod)

	for {
		select {
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(60))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}

			var data DiffDepth
			err = json.Unmarshal(msg, &data)
			if err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
			c <- data
		}
	}
	return nil
}

type DiffDepth struct {
	E  string     `json:"e"`
	E0 int        `json:"E"`
	S  string     `json:"s"`
	U  int        `json:"U"`
	U0 int        `json:"u"`
	B  [][]string `json:"b"`
	A  [][]string `json:"a"`
}
