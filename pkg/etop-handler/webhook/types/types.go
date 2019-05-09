package types

// MessageCollector is a non-thread safe collector which manages messages and size
type MessageCollector struct {
	Messages [][]byte
	Size     int
}

func (c *MessageCollector) Collect(msg []byte) {
	c.Messages = append(c.Messages, msg)
	c.Size += len(msg)
}
