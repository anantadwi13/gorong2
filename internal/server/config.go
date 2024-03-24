package server

type Config struct {
	Backbone ConfigBackbone `json:"backbone"`
	Auth     ConfigAuth     `json:"auth"`
}

type ConfigBackbone struct {
	ListenAddr string     `json:"listen_addr"` // "localhost:19000"
	ListenType ListenType `json:"listen_type"` // "websocket", "tcp"

	TLSEnable bool   `json:"tls_enable"`
	TLSCert   string `json:"tls_cert"` // filepath
	TLSKey    string `json:"tls_key"`  // filepath
	TLSCa     string `json:"tls_ca"`   // filepath
}

type ConfigAuth struct {
	AuthType AuthType `json:"auth_type"`

	PlainToken string `json:"plain_token"`
}

type (
	ListenType string
	AuthType   string
)

const (
	ListenTypeWebsocket ListenType = "websocket"
	ListenTypeTcp       ListenType = "tcp"

	AuthTypePlainToken AuthType = "plain-token"
)
