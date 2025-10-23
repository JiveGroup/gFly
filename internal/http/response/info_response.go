package response

// ServerInfo struct to describe system information.
// @Description contains system metadata including name, server prefix, and server name.
// @Name Name is the name of the API.
// @Prefix Prefix is the API prefix including version.
// @Server Server is the name of the server application.
// @Tags Success Responses
type ServerInfo struct {
	Name   string `json:"name" example:"ThietNgon API" doc:"API name"`
	Prefix string `json:"prefix" example:"/api/v1" doc:"API prefix including version"`
	Server string `json:"server" example:"ThietNgon-Go Server" doc:"Server application name"`
}
