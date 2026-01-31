package dto

type GoogleAuthRequest struct {
	IDToken string `json:"idToken" binding:"required"`
}

type GoogleAuthResponse struct {
	Token string            `json:"token"`
	User  GoogleUserDetails `json:"user"`
}

type GoogleUserDetails struct {
	Email     string `json:"email"`
	Fullname  string `json:"fullname"`
	IsNewUser bool   `json:"isNewUser"`
}
