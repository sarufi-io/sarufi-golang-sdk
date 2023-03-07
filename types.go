package sarufi

// Application will hold all application
// methods.
type Application struct {
}

// Type Bot. All the fields are matched to the API JSON response.
// You can read more here https://neurotech-africa.stoplight.io/docs/sarufi/a3135fbb09470-create-new-chatbot
type Bot struct {
	Id                    int                    `json:"id"`
	Name                  string                 `json:"name"`
	Industry              string                 `json:"industry"`
	Description           string                 `json:"description"`
	UserID                int                    `json:"user_id"`
	Username              string                 `json:"user_name"`
	VisibleOnCommunity    bool                   `json:"visible_on_community"`
	Intents               map[string][]string    `json:"intents"`
	Flows                 map[string]interface{} `json:"flows"`
	ModelName             string                 `json:"model_name"`
	WebhookURL            string                 `json:"webhook_url"`
	WebhookTriggerIntents []string               `json:"webhook_trigger_intents"`
	EvaluationMetrics     interface{}            `json:"evaluation_metrics"`
	CreatedAt             string                 `json:"created_at"`
	UpdatedAt             string                 `json:"updated_at"`
	ChatID                string                 `json:"chat_id"`
	Conversation          Conversation
	Prediction            Prediction
	ChatUsers             []ChatUser
	ConversationHistory   []ConversationHistory `json:"conversation_history"`
}

// This type will be used in the conversation history
type Response struct {
	SendMessage []string `json:"send_message"`
}

// This type will hold information of the converstion
// of a particular chat ID
type ConversationHistory struct {
	ID           int        `json:"id"`
	Message      string     `json:"message"`
	Sender       string     `json:"sender"`
	Response     []Response `json:"response"`
	ReceivedTime string     `json:"received_time"`
}

// ChatUser type to hold information about
// a chat
type ChatUser struct {
	ChatID       string `json:"chat_id"`
	ReceivedTime string `json:"received_time"`
}

// Prediction will hold the response
// of a paritcular intent prediction
type Prediction struct {
	Message    string  `json:"message"`
	Intent     string  `json:"intent"`
	Status     bool    `json:"status"`
	Confidence float64 `json:"confidence"`
}

// Conversation will hold all conversation
// related data
type Conversation struct {
	Message      []string          `json:"message"`
	Memory       map[string]string `json:"memory"`
	CurrentState string            `json:"current_state"`
	NextState    string            `json:"next_state"`
}

// Type Flow to create a new Flow XD.
// Message is a slice of strings.
// NextState is simply a string.
type Flow struct {
	Message   []string `json:"message"`
	NextState string   `json:"next_state"`
}

// Type Flows to map a string to the type Flow.
type Flows map[string]interface{}
