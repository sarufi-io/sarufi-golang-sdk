// Golang SDK for Sarufi Conversational AI Platform
package sarufi

// GetToken() method to generate a token for the user.
// The received token will be saved to Application.Token
// otherwise it will return an error.
func (app *Application) SetToken(apiKey string) {
	token = apiKey
}
