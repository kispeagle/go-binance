package stream

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	logger "github.com/kispeagle/go-binance/logs"
)

type TickerService struct {
	Symbol string
}

func NewTickerService(symbol string) TickerService {
	s := TickerService{Symbol: symbol}
	return s
}

func (s TickerService) Call(url string, c chan interface{}) error {

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	conn.SetReadLimit(4096)
	ticker := time.NewTicker(PingPeriod)
	if err != nil {
		logger.Log.Debug(err)
		c <- nil
		return err
	}

	for {
		select {
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(60))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c <- nil
				c <- nil
				return err
			}
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				c <- nil
				c <- nil
				return err
			}

			var data Ticker
			json.Unmarshal(msg, &data)
			c <- nil
			c <- data
		}
	}
	return nil
}

type Ticker struct {
	E  string `json:"e"`
	E0 int    `json:"E"`
	S  string `json:"s"`
	P  string `json:"p"`
	P0 string `json:"P"`
	W  string `json:"w"`
	X  string `json:"x"`
	C  string `json:"c"`
	Q  string `json:"Q"`
	B  string `json:"b"`
	B0 string `json:"B"`
	A  string `json:"a"`
	A0 string `json:"A"`
	O  string `json:"o"`
	H  string `json:"h"`
	L  string `json:"l"`
	V  string `json:"v"`
	Q0 string `json:"q"`
	O0 int    `json:"O"`
	C0 int    `json:"C"`
	F  int    `json:"F"`
	L0 int    `json:"L"`
	N  int    `json:"n"`
}
