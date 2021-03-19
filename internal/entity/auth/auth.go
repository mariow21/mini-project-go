package auth

// UserID ...
type UserID struct {
	UserID int `json:"userid"`
}

// Metadata ...
type Metadata struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

// Err ...
type Err struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Code   int    `json:"code"`
}

// Auth ...
type Auth struct {
	Data     interface{} `json:"data"`
	Metadata Metadata    `json:"metadata"`
	Error    Err         `json:"error"`
}
