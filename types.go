package sarufi

// Type Bot. All the fields are matched to the API JSON response.
// You can read more here https://neurotech-africa.stoplight.io/docs/sarufi/a3135fbb09470-create-new-chatbot
type Bot struct {
	Id                        int                    `json:"id"`
	Name                      string                 `json:"name"`
	Industry                  string                 `json:"industry"`
	Description               string                 `json:"description"`
	UserID                    int                    `json:"user_id"`
	VisibleOnCommunity        bool                   `json:"visible_on_community"`
	Intents                   map[string][]string    `json:"intents"`
	Flows                     map[string]interface{} `json:"flows"`
	ModelName                 string                 `json:"model_name"`
	WebhookURL                string                 `json:"webhook_url"`
	WebhookTriggerIntents     []string               `json:"webhook_trigger_intents"`
	EvaluationMetrics         interface{}            `json:"evaluation_metrics"`
	ChatID                    string                 `json:"chat_id"`
	Conversation              Conversation
	ConversationWithKnowledge ConversationWithKnowledge
	Prediction                Prediction
	ChatUsers                 []ChatUser
	ConversationHistory       []ConversationHistory `json:"conversation_history"`
}

// This type will be used in the conversation history
type Response struct {
	Message []string `json:"send_message"`
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
// related data on none knowledge base bot
type Conversation struct {
	Message   []string `json:"message"`
	Memory    Memory   `json:"memory"`
	NextState string   `json:"next_state"`
}

// ConversationWithKnowledge will hold all conversation
// related data on knowledge base bot
type ConversationWithKnowledge struct {
	Message   []Actions `json:"actions"`
	Memory    Memory    `json:"memory"`
	NextState string    `json:"next_state"`
}

type Actions struct {
	ResponseMessage []any `json:"send_message"`
}

type Memory interface{}

// Application will hold all application methods.
// It also has details about the user
type Application struct {
	User User `json:"user"`
}

type User struct {
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Mobile      string `json:"mobile"`
	IsAdmin     bool   `json:"is_admin"`
	DateCreated string `json:"date_created"`
	UpdatedAt   string `json:"updated_at"`
}
