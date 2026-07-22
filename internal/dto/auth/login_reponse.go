package dto

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"-"` //This prevents the refresh token from ever being serialized into JSON.
}
