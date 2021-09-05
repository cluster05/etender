package auth

type JwtMeta struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	ExpiredAt int64  `json:"expiredAt"`
}

type AccessToken struct {
}
