package response

// SignIn struct to describe login response.
type SignIn struct {
	Access  string `json:"access" doc:"The access token for authentication"`
	Refresh string `json:"refresh" doc:"The refresh token for obtaining a new access token"`
}
