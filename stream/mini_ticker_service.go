package stream

import (
	"encoding/json"
	"time"

	logger "github.com/kispeagle/go-binance/logs"

	"github.com/gorilla/websocket"
)

type MiniTickerService struct {
	Symbol string
}

func NewMiniTickerService(symbol string) MiniTickerService {
	miniTicker := MiniTickerService{symbol}
	return miniTicker
}

func (s MiniTickerService) Call(url string, c chan interface{}) error {

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	ticker := time.NewTicker(PingPeriod)

	if err != nil {
		logger.Log.Error(err)
		c <- nil
		return err
	}

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
			var data MiniTicker
			err = json.Unmarshal(msg, &data)
			if err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
			c <- data
		}
	}
}

type MiniTicker struct {
	E  string `json:"e"`
	E0 int    `json:"E"`
	S  string `json:"s"`
	C  string `json:"c"`
	O  string `json:"o"`
	H  string `json:"h"`
	L  string `json:"l"`
	V  string `json:"v"`
	Q  string `json:"q"`
}
