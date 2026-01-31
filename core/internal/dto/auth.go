package dto

type LoginBody struct {
	IDToken string  `json:"idToken" binding:"required"`
	Name    *string `json:"name"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
