package stream

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	logger "github.com/kispeagle/go-binance/logs"
)

type AllMiniTickersService struct {
	Symbol string
}

func NewAllMiniTickersService(symbol string) AllMiniTickersService {
	s := AllMiniTickersService{symbol}
	return s
}

func (s AllMiniTickersService) Call(url string, c chan interface{}) error {

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

			var data MiniTickerList
			err = json.Unmarshal(msg, &data.List)
			if err != nil {
				logger.Log.Error(err)
				c <- nil
				return err
			}
			c <- data
		}
	}

}

type MiniTickerList struct {
	List *[]MiniTicker
}
