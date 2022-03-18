package main

type Heartbeat struct {
}

func NewHeartbeat() *Heartbeat {
	return &Heartbeat{}
}

func (h *Heartbeat) GetType() string {
	return "heartbeat"
}

func (h *Heartbeat) GetId() string {
	return ""
}
