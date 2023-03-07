package main

import (
	"fmt"
	"log"

	"github.com/sarufi-io/sarufi-golang-sdk"
)

func main() {
	var app sarufi.Application

	// Getting Token
	if err := app.GetToken("string", "string"); err != nil {
		log.Fatal(err)
	}

	// Getting all bots
	myBots, err := app.GetAllBots()
	if err != nil {
		log.Fatal(err)
	}

	for _, bot := range myBots {
		fmt.Printf("%d: %s\n", bot.Id, bot.Name)
	}

	// Getting a single bot
	example_bot, err := app.GetBot(myBots[0].Id)

	if err != nil {
		log.Fatal(err)
	}

	// Updating a bot
	example_bot.Name = "New Name"

	if err = app.UpdateBot(example_bot); err != nil {
		log.Fatal(err)
	}

	// Deleting a bot
	if err = app.DeleteBot(449); err != nil {
		log.Fatal(err)
	}

	// Create a new bot
	example_bot, err = app.CreateBot("Name of your bot", "Description", "Industry", false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(example_bot.Name)

	// Creating new intents
	intents := `
	{
	    "goodbye": ["bye", "goodbye", "see ya"],
	    "greets": ["hey", "hello", "hi"],
	    "order_pizza": ["I need pizza", "I want pizza"]
	}`

	if err := example_bot.CreateIntents(intents); err != nil {
		log.Fatal(err)
	}

	// For changes to take effect
	if err = app.UpdateBot(example_bot); err != nil {
		log.Fatal(err)
	}

	// create a string in JSON format
	newFlow := `
	{
	     "greets": {"message": ["Hi, How can I help you?"], "next_state": "end"},
	     "order_pizza": {
	     "message": ["Sure, How many pizzas would you like to order?"],
	     "next_state": "number_of_pizzas"
	 },
	     "number_of_pizzas": {
	     "message": [
	         "Sure, What would you like to have on your pizza?",
	         "1. Cheese",
	         "2. Pepperoni",
	         "3. Both"
	     ],
	     "next_state": "choice_pizza_toppings"
	 },
	     "choice_pizza_toppings": {
	     "1": "pizza_toppings",
	     "2": "pizza_toppings",
	     "3": "pizza_toppings",
	     "fallback_message": ["Sorry, the topping you chose is not available."]
	 },
	     "pizza_toppings": {
	     "message": ["Cool, Whats your address ?"],
	     "next_state": "address"
	 },
	     "address": {
	     "message": ["Sure, What is your phone number ?"],
	     "next_state": "phone_number"
	 },
	     "phone_number": {
	     "message": ["Your order has been placed.", "Thank you for ordering with us."],
	     "next_state": "end"
	 },
	     "goodbye": {"message": ["Bye", "See you soon"], "next_state": "end"}
	}`

	// I am ignoring errors but you should handle them
	example_bot.CreateFlows(newFlow)
	app.UpdateBot(example_bot)

	fmt.Println(example_bot)

	// Respond to a message
	if err = example_bot.Respond("Hey", "general"); err != nil {
		log.Fatal(err)
	}

	fmt.Println(example_bot.Conversation.Message)

	// Get chat state
	if err := example_bot.ChatState(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(example_bot.Conversation.CurrentState)

	// Get Chat History
	if err := example_bot.GetChatHistory(example_bot.ChatID); err != nil {
		log.Fatal(err)
	}

	for _, chat := range example_bot.ConversationHistory {
		fmt.Printf("id: %d\nmessage: %s\nsender: %s\nresponse: %v\nreceived time: %s\n\n", chat.ID, chat.Message, chat.Sender, chat.Response, chat.ReceivedTime)
	}

	// Get Chat Users
	if err := example_bot.GetChatUsers(); err != nil {
		log.Fatal(err)
	}

	for _, chat := range example_bot.ChatUsers {
		fmt.Println(chat.ChatID)
		fmt.Println(chat.ReceivedTime)
	}

	// Predict a message confidence
	if err := example_bot.Predict("Hey"); err != nil {
		log.Fatal(err)
	}

	fmt.Println(example_bot.Prediction.Confidence)

	// Get User information

	if err := app.GetUser(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(app.User.ID)
	fmt.Println(app.User.Username)
}
