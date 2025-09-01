package server

type SetKeyRequest struct {
	Key             string `json:"key"`
	Value           string `json:"value"`
	ExpiresAtMillis uint64 `json:"expiresAt"`
}

type DeleteKeyRequest struct {
	Key string `json:"key"`
}

type GetKeyRequest struct {
	Key string `json:"key"`
}
