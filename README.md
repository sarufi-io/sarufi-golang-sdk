<h1 align="center">Sarufi Golang SDK</h1>
<p align="center">
<img src="https://img.shields.io/badge/Made%20with-Go-1f425f.svg">
<img src="https://img.shields.io/github/go-mod/go-version/gomods/athens.svg">
<a href="https://pkg.go.dev/github.com/sarufi-io/sarufi-golang-sdk"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>
</p>
Golang SDK for [Sarufi Conversational AI Platform](https://docs.sarufi.io/). Made with love for Gophers.

## Pre-requisites
You need to register with [Sarufi](https://docs.sarufi.io/) in order to get your login credentials.

## Installation
Simply install with the `go get` command:
```go
go get github.com/sarufi-io/sarufi-golang-sdk@v0.2.0
```

And import it to your package as:
```go
package main

import (
    sarufi "github.com/sarufi-io/sarufi-golang-sdk@v0.2.0"
)
```

## Usage
I have divided the explanation into two parts;
- Information relating to type `Sarufi.Application` and 
- Information relating to type `Sarufi.Bot`

## Type Sarufi.Application 
This type handles tasks that affect all bots such getting all bots, creating new bots and so on as explained below.

### Getting Token 
To get a token, first create a new variable of type `sarufi.Application` and sign in with your credentials as shown below. The method will return an error if any:
```go
var app sarufi.Application

if err := app.GetToken("username", "password"); err != nil {
    log.Fatal(err)
    }
```

### Creating a New Bot
Use the `app.CreateBot` method to create a new bot. You'll fill in the name of your bot, its description, industry and a boolean on whether the bot should be visible on the playground.
The method will return a new bot of type `Bot` and an error if any. See below:

```go
example_bot, err := app.CreateBot("Name of your bot", "Description", "Industry", bool)
if err != nil {
        log.Fatal(err)
    }
fmt.Println(example_bot.Name)
```

### Getting All Your Bots
To get all your bots, use the `app.GetAllBots` which will return a list of all available bots of type `Bot` and an error if any:
```go
myBots, err := app.GetAllBots()
if err != nil {
        log.Fatal(err)
    }

for _, bot := range myBots {
        fmt.Println(bot.Name)
        fmt.Println(bot.Description)
    }
```

### Fetch A Single Bot
To fetch a particular bot, use the method `app.GetBot` filling in the ID of the bot you want. It will return the requested bot and error if any:
```go
bot, err := app.GetBot(bot_id)
if err != nil {
    log.Fatal(err)
    }

fmt.Println(bot.Name)
```

### Updating A Bot
To update information of a specific bot, use the `app.UpdateBot` placing the bot itself as a parameter. It will return an error if any:
```go
if err := app.Update(example_bot); if err != nil {
        log.Fatal(err)
    }
```

### Delete A Bot
To delete a bot, use the `app.DeleteBot` method filling in the ID of the bot as a parameter. It will return an error if any:
```go
if err := app.Delete(bot_id); if err != nil {
        log.Fatal(err)
    }
```

## Type Sarufi.Bot 
This is the actual `Bot`. Below are explanations on all methods that affect a specific bot such as creating new intents, flows and so on.

**NOTE: For any changes to take effect, you MUST call the app.Update method with the bot as a parameter, otherwise changes will not take effect. More explanation below.**

### Creating A New Intent
You can create a new Intent with the `bot.CreateIntents` method. Note that this WILL delete any existing intents. 
```go
// Create a string in JSON format

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
```

- You can also create a new `Intent` with a JSON file as shown below:
```go
example_bot.CreateIntents("intent.json")
```

### Adding An Intent

You can add a new `Intent` to the ones already available as follows:
```go
// Create a title of the new intent

title := "example_title"

// Create a slice of strings with intent messages

newIntent := []string{"example intent content"}

// I am ignoring errors but you should handle them 
example_bot.AddIntent(title, newIntent)
app.Update(example_bot)
```
*Note that this WILL NOT delete previous `Intents`. It will simple add itself to them.*

### Deleting An Intent
You can delete a specific `Intent` as follow:
```go
// I am ignoring errors but you should handle them 
example_bot.DeleteIntent("intent_title")
app.UpdateBot(example_bot)
```

### Creating A New Flow
You can create a new `Flow` as follows:
```go
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
```

- You can also create a new `Flow` with a JSON file as shown below:
```go
example_bot.CreateFlows("flow.json")
```


### Adding A Flow
Add a new `Flow` to the already existing ones as follows:
```go
// Create a title string

title := "flow_title"
content := `{
"message": ["Your order has been placed.", "Thank you for ordering with us."],
"next_state": "end"
}`

// I am ignoring errors but you should handle them 
example_bot.AddFlow(title, content)
app.UpdateBot(example_bot)
```
*Note that this will NOT delete previous `Flows`. It will simple add itself to them.*

### Deleting A Flow
You can delete a specific `Flow` as follow:
```go
// I am ignoring errors but you should handle them 
example_bot.DeleteFlow("flow_title")
app.UpdateBot(example_bot)
```

### Respond To A Message
Use the `bot.Respond` method to get a bot's respond to a particular message. The response will be stored at `bot.Conversation` which has the fields message, memory, current and next states.
```go
if err = example_bot.Respond("Hey", "general"); err != nil {
		log.Fatal(err)
	}

fmt.Println(example_bot.Conversation.Message)
```

### Check Chat State
You can check the current and the next state of the chat using the `bot.ChatState` method. The states are stored at the `bot.Conversation` field.
```go
if err := example_bot.ChatStatus(); err != nil {
        log.Fatal(err)
    }

fmt.Println(example_bot.Conversation.CurrentState)

```

### Get Chat History
You can fetch the history of a specific chat ID using the `bot.GetChatHistory` method passing in the chat ID as a parameter. The history will be saved in the `bot.ConversationHistory` field.
```go
if err := bot.GetChatHistory("chat_id"); err != nil {
        log.Fatal(err)
    }

for _, chat := range bot.ConversationHistory {
	fmt.Printf("id: %d\nmessage: %s\nsender: %s\nresponse: %v\nreceived time: %s\n\n", chat.ID, chat.Message, chat.Sender, chat.Response, chat.ReceivedTime)
	}
```

### Get Chat Users
To get a list of all users communicating with your bot, use the `bot.GetChatUsers` method. The list of chats information will be stored in the `bot.ChatUsers` field.
```go
if err := bot.GetChatUsers(); err != nil {
        log.Fatal(err)
    }

for _, chat := range bot.ChatUsers {
        fmt.Println(chat.ChatID)
        fmt.Println(chat.ReceivedTime)
    }
```

### Predict A Message 
To get a prediction of a particular message on your bot, use the `bot.Predict` method with the message as a parameter. The result of the prediction will be stored in the `bot.Prediction` field.
```go
if err := bot.Predict("Hey"); err != nil {
        log.Fatal(err)
    }

fmt.Println(bot.Prediction.Confidence)
```

## Additional Resources
- https://docs.sarufi.io/
- https://neurotech-africa.stoplight.io/docs/sarufi 

## Issues
Open new issues if any.

## Thanks
- [Pius](https://github.com/piusalfred)
- [Kalebu](https://github.com/Kalebu)
- All other [contributors](https://github.com/sarufi-io/sarufi-golang-sdk/graphs/contributors)

## Author
This package is authored and maintained by [Avicenna](https://github.com/AvicennaJr)
