# instaspy

## Project description
- [Description](#description)
- [Installation](#installation)
- [Get started](#Get-started)
- [Autorun](#autorun)
- [ToDo's](#todo)

## Description

**Act like you don't see me, we'll play pretend**  

This tool is litteraly your ex's best friend to spy on you. Why? This tool will help you to spy on a people instagram stories. Once new image have been added you will be informed with telegram message. You will recieve image itself and string telling whom this picture is posted by(inst: *username*). Pretty minimalist huh?

## Installation 
Before we get into installation process make sure that you have docker and docker compose installed on your system.

To install application it's enough to clone repo. Assuming that you will be running this application on server so you could spy on someone 24/7 make those steps there.
```bash
git clone https://github.com/panaglev/instaspy.git
cd instaspy
```

## Get started
First of all we need to edit environment variables. Rename .env.template to .env and open it with any text editor that you like.
```bash
mv .env.template .env
```
You will see there two empty fields that you need to fill. TELEGRAM_BOT is telegram bot token and CHAT_ID is the field with chat id where your bot acting. Change those fields and save .env file, ex:
```bash
TELEGRAM_BOT=123123123:123123123:123123123
CHAT_ID=123123123
```
Later we need to specify which instagram account we want to spy on. Open config.yaml file and edit first line with usernames, ex:
```yaml
usernames: ["miakhalifa", "jialissaonly", "letrileylive"]
```
After configuration is set up change start_and_watch.sh script permission and feel free to run app.
```bash
chmod +x start_and_watch.sh
```
Not really sure if docker compose gonna create volume automatic so please don't be lazy to write:
```bash
docker volume create mydata
```
This volume is used for storing db with info about pictures that have been already sent.

## Autorun
If you want app to run automatically you might use crontab to add it there. For example I use it every 10 minutes. To add job to crontab use:
```bash
crontab -e
```
After you have opened crontab file to edit add line below in the end of file:
```bash
*/10 * * * * cd /root/instaspy && ./start_and_watch.sh
```

## Todo
- Add video parce(for now it's only in pictures mode)
- Using for parces purposes self-written service instead of parsing other sites
- To finally find a job...

Let me know if you want me to become your employee(ru/en) 