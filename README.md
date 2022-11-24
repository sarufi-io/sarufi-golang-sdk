# Sarufi Golang SDK

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/sarufi-io/sarufi-golang-sdk)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/sarufi-io/sarufi-golang-sdk)

Golang SDK for [Sarufi Conversational AI Platform](https://docs.sarufi.io/). Made with love for Gophers.

## Pre-requisites

You need to register with [Sarufi](https://docs.sarufi.io/) in order to get your login credentials.

## Installation

Simply install with the `go get` command:
```
go get github.com/sarufi-io/sarufi-golang-sdk
```
And import it to your package as:
```
package main

import (
    sarufi "github.com/sarufi-io/sarufi-golang-sdk"
)
```

## Usage

First create a new variable of type `sarufi.Bot` and initialze it as shown below:
```
var bot sarufi.Bot
bot.Initialize("your-email", "your-password")
```

### Creating a New Bot

Simply create a new bot with:

```
bot.Create("Name of your bot", "Description", "Industry")
```

### Getting All Your Bots

Get a JSON response of all your bots with the `GetBots()` method:
```
bot.GetBots()
```

### Fetch A Single Bot

You can select a particular bot with the `GetBot()` method passing the bot id as a parameter:
```
bot.GetBot(100)

// It will print out your bots name

fmt.Println(bot.Name)
```

### Creating A New Intent

You can create a new Intent as follows:
```
// Create a string in JSON format

intents := `
{
    "goodbye": ["bye", "goodbye", "see ya"],
    "greets": ["hey", "hello", "hi"],
    "order_pizza": ["I need pizza", "I want pizza"]
}`

bot.CreateIntents(intents)
```

You can also create a new `Intent` with a JSON file as shown below:
```
bot.CreateIntents("intent.json")
```

### Adding An Intent

You can add a new `Intent` to the ones already available as follows:
```
// Create a title of the new intent

title := "football"

// Create a slice of strings with intent messages

newIntent := []string{"What team do you support?"}

// Add them to the bot as follows:

bot.AddIntent(title, newIntent)

// The bot will now have a new intent call football, in addition to the ones that were already available
```
Note that this will NOT delete previous `Intents`. It will simple add itself to them.

### Deleting An Intent

You can delete a specific `Intent` as follow:
```
bot.DeleteIntent("Intent Title")
```

### Creating A New Flow

You can create a new `Flow` as follows:
```
// create a string in JSON format

newFlow := `
{
	"football": {
		"message": ["Do you like football"],
		"next_state": "team"
	},
	"goodbye": {
		"message": ["Bye", "See you soon"],
		"next_state": "end"
	},
	"greets": {
		"message": ["Hi, How can I help you?"],
		"next_state": "end"
	}
}`

bot.CreateFlows(newFlow)
```

You can also create a new `Flow` with a JSON file as shown below:
```
bot.CreateFlows("flow.json")
```


### Adding A Flow

Add a new `Flow` to the already existing ones as follows:
```
// Create a title string

title := "flowabanga"

// Create a new variable of type Flow

var newFlow bot.Flow

// Fill in the new messages and the next state

newFlow.Message = []string{"A new list of messages"}
newFlow.NextState = "To the next state"

bot.AddFlow(title, newFlow)
```
Note that this will NOT delete previous `Flows`. It will simple add itself to them.

### Deleting A Flow

You can delete a specific `Flow` as follow:
```
bot.DeleteFlow("Flow Title")
```

### Respond To A Message

You can test how a bot will respond to a message as follows:
```
bot.Respond("order pizza", "Text") 

// Message followed by the message type. The message type here is text.
```

### Check Chat State

You can check the current and the next state of the chat as follows:
```
bot.ChatState()
```

### Delete A Bot

You can delete a bot as follows:
```
bot.Delete(100) // Pass the bot ID
```

## Additional Resources

To learn more about Sarufi, kindly visit https://docs.sarufi.io/.

## Issues

The Golang SDK does NOT support `Flows` with choices yet.

## Author

This package is authored and maintained by [Mojo](https://github.com/AvicennaJr)
