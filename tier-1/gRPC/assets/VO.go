package assets

type CryptoRequestJSON struct {
	Hash         string `json:"hash"`
	SignatureB64 string `json:"signature_b64"`
}

type CryptoResponseJSON struct {
	Valid   bool   `json:"valid"`
	Message string `json:"menssage"`
}
