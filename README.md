# Jester Bot

## Features

* Stable Diffusion AI Integration
* ITZG Minecraft Server Generator 

# Setup

* Register a discord bot
* export BOT_TOKEN to your bot's token
* Enable all Privileged Gateway Intents in discord

* Install https://github.com/AUTOMATIC1111/stable-diffusion-webui
* Enable API by editing `webui-user.bat` to include `set COMMANDLINE_ARGS=--xformers --autolaunch --listen --api`
* Start server

# Usage

* Run bot with `go run main.go -t $BOT_TOKEN`
* Send command !img followed by a prompt for the txt2img AI
* Send command !server followed by ["start", "stop", "save"]

# Build

# Tests

# Commands

* !img a sunny winter day in colorado