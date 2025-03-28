# docker build --tag tbot-chatgpt:0.0.1 .
# docker run --restart unless-stopped --name tbot-chatgpt -e TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN -e CHATGPT_API_TOKEN=$CHATGPT_API_TOKEN -e CHATGPT_API_DB_URI=$CHATGPT_API_DB_URI -e CHATGPT_API_BOLT_FILE=$CHATGPT_API_BOLT_FILE -d tbot-chatgpt:0.0.1
FROM alpine:3.20.0
WORKDIR /
COPY ./bin/tbot-chatgpt-amd64 .

CMD ["/tbot-chatgpt-amd64"]