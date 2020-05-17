# FeedThemBot
A telegram bot to remind user when to eat during the day.  
This bot helps to keep up with a diet or just eating more regularly.
### Bot creation
To create a bot one should write to BotFather.  
Once the bot is created, one will receive an authorization token which is required to start using the bot.
### Deploy and Run the bot in Docker
Before running this bot on their machine, one should add their bot token to `docker-compose.yml` as an environment variable.
```
 feedthembot:
   ...
   environment:
     TOKEN: Your-Token-Right-Here
```
After that, deploy and run by executing the following command in the project main directory
```
$ docker-compose up --build
```
This program would deploy database to separate Docker container.  
In order to use another database simply specify its data in `docker-compose.yml`
```
 environment:
   ...
   environment:
     PGHOST: <db_hostname>
     PGDATABASE: <db_name>
     ...
```
After that, execute the following command
```
$ docker-compose up --build feedthembot
```
