package terminal

const (
	WsMsgCmd       = "cmd"
	WsMsgResize    = "resize"
	WsMsgHeartbeat = "heartbeat"
	// WsMsgSession is the first message after attach; carries the session id.
	WsMsgSession = "session"
)

type WsMsg struct {
	Type      string `json:"type"`
	Data      string `json:"data,omitempty"`     // WsMsgCmd
	Cols      int    `json:"cols,omitzero"`      // WsMsgResize
	Rows      int    `json:"rows,omitzero"`      // WsMsgResize
	Timestamp int    `json:"timestamp,omitzero"` // WsMsgHeartbeat
	ID        string `json:"id,omitempty"`       // WsMsgSession
}

func setQuit(ch chan bool) {
	ch <- true
}
