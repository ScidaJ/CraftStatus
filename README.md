# CraftStatus

This is a small self hosted Discord bot designed to monitor a Minecraft server. It features player count monitoring, server status in the sidebar, and structured logging.

## Running the Bot

If you're just looking for a `docker-compose.yml` or `.exe` to run then head on over to the releases. You'll find ZIPs of both, as well as an `.env.sample` and a README that is a copy of this one.

Here is a link to the Image on Docker Hub: https://hub.docker.com/repository/docker/scidaj57/minecraft-helper/general

### Docker Compose

I'm not sure the best way to do this so this is my `docker-compose.yml` that I am currently using.

```YAML
services:
  bot:
    container_name: bot
    image: "scidaj57/minecraft-helper"
    ports:
      - "8080:8080"
    env_file:
      - "./.env"
```

You can also define your environment variables within the `docker-compose.yml` like so

```YAML
services:
  bot:
    container_name: bot
    image: "scidaj57/minecraft-helper"
    ports:
      - "8080:8080"
    environment:
      BOT_TOKEN:      "XXX0XXXxXxX0XXXxXxXxXXX0Xx"
      GUILD_ID:       "000000000000000000"
      RCON_ADDRESS:   "127.0.0.1:25565"
      RCON_PASSWORD:  "hunter2"
      ADMIN:          "000000000000000000" 
```

## Discord Setup

This requires making an application with Discord on the Discord Developer Portal, found [here.](https://discord.com/developers/applications) There are a few pieces of information that we need from there, so lets go over what they are.

### Bot Token

Once you have created your application for the bot to work with, look on the sidebar for the `Bot` category, as that is where we'll find our `Bot Token`. Once on the Bot page, click the `Reset Token` button under the `Username` field. Copy this token and place it in your `.env` file.

### Permissions

After the previous step head over to the `OAuth2` category on the sidebar, then scroll down to the `OAuth2 URL Generator`. Here you're going to want to select the `bot` scope, in the middle column. A second `Bot Permissions` panel will open underneath, where we will define what permissions the bot needs. As of now, it only needs the following permissions.

* General Permissions
  * Read Messages/View Channels
* Text Permissions
  * Send Messages
  * Manage Messages
  * Read Message History

Select these in the `Bot Permissions` panel, then copy the `Generated URL` below and paste it in a new tab. From here select the server that you wish to add the bot to, and confirm. When you launch the bot with the associated server ID your application should appear in the sidebar.

<a id=".env"></a>
### .env

This will be a quick overview of the variables in the `.env` file.

* `BOT_TOKEN` - The token the bot will use to log in to Discord with. If it does not match your bot token from your application in the Discord Developer Portal then the request will be denied.
* `GUILD_ID` - The ID of the Discord server that you would like the bot to join upon start up. The bot must be added to the server first with the required permissions.
* `RCON_ADDRESS` - This is set in your `server.properties` file or similar. Port must be supplied with the address.
* `RCON_PASSWORD` - This is set in your `server.properties` file or similar.
* `ADMIN` - The User ID of the "Admin" user for the bot/server. They will be pinged if there is an issue with the server.
* `SERVER_ADDRESS` Optional. The `/address` command just returns the IP of the host machine, as this bot is assuming that the server and bot are running on the same machine. If this variable is filled in then it will instead return this value.
* `SERVER_PORT` Optional. Will append to the end of the `SERVER_ADDRESS` value if present.
* `PLAYER_LIST` Optional. For use with `/list` command. If value is provided in the format of `[InGameName1:Nickname1,InGameName2:Nickname2,InGameName3:Nickname3]` then it will replace the in game name with the provided nickname in the list. If no nickname is provided then it will print the in game name instead.

## Commands

The bot only has three commands as it is fairly simple in scope. They are listed below

* `/address` Prints out the address of the server. As the bot assumes that the server is running on the same machine it will return the IP of the host machine. **If you do not want this to be the case then fill in the `SERVER_ADDRESS` variable in the `.env` file. It will print that value instead.**
* `/list` List the players currently on the server. If the `.env` variable `PLAYER_LIST` is populated the bot will replace any matching usernames with the corresponding nickname.
* `/restart` Sends the `/stop` command to the server then checks every 30 seconds for 5 minutes if it has relaunched. You must have some way of restarting your server automatically. 
