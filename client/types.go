package client

type MessageItems struct {
	Type int
	Content []byte
}


const (
	CoinexDepth = iota + 1
)