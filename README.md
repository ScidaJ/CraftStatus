# DiscordMinecraftHelper

This is a small self hosted Discord bot designed to monitor a Minecraft server. It features player count monitoring, server status in the sidebar, structured logging, and automatic daily restarts at 2AM. Currently only Windows compatible.

## Requirements

`go >= 1.21`

To use the daily restarts and `/restart-server` command, you need a bat file which automatically restarts the server upon shutdown. The All The Mods 9 modpack has one, but all you need is to add this snippet to your `startserver.bat`, around the existing code that launches the `server.jar`

```bat
:START
:: Code that launches server jar
echo Restarting automatically in 10 seconds (press Ctrl + C to cancel)
timeout /t 10 /nobreak > NUL
goto:START
```
and you're good to go!

## Installing

1. Clone this repo ```git clone git@github.com:ScidaJ/DiscordMinecraftHelper.git```
2. CD into the new directory ```cd DiscordMinecraftHelper```
3. Install dependencies ```go mod download```
4. Make a copy of `.env.sample` and fill in the required values
5. Run the bot with ```go run main.go```

## Everything Else

This requires making an application with Discord on the Discord Developer Portal, found [here.](https://discord.com/developers/applications) There are a few pieces of information that we need from there, so lets go over what they are.

### Bot Token

Once you have created your application for the bot to work with, look on the sidebar for the `Bot` category, as that is where we'll find our `Bot Token`. Once on the Bot page, click the `Reset Token` button under the `Username` field. Copy this token and place it in your .env file.

### Permissions

After the previous step head over to the `OAuth2` category on the sidebar, then scroll down to the `OAuth2 URL Generator`. Here you're going to want to select the `bot` scope, in the middle column. A second `Bot Permissions` panel will open underneath, where we will define what permissions the bot needs. As of now, it only needs the following permissions.

* General Permissions
  * Read Messages/View Channels
* Text Permissions
  * Send Messages
  * Manage Messages
  * Read Message History

Select these in the `Bot Permissions` panel, then copy the `Generated URL` below and paste it in a new tab. From here select the server that you wish to add the bot to, and confirm. When you launch the bot with the associated server ID your application should appear in the sidebar.
