package stream

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	logger "github.com/kispeagle/go-binance/logs"
)

type BookTickerService struct {
	Symbol string
}

func NewBookTickerService(symbol string) BookTickerService {
	s := BookTickerService{symbol}
	return s
}

func (s BookTickerService) Call(url string, c chan interface{}) error {

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

			var data BookTicker
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

type BookTicker struct {
	U  int    `json:"u"`
	S  string `json:"s"`
	B  string `json:"b"`
	B0 string `json:"B"`
	A  string `json:"a"`
	A0 string `json:"A"`
}
