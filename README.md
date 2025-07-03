# Insider Assessment Project

## Description
Automatic message sending system:
- Every 2 minutes, fetch 2 unsent messages from the DB, send via webhook, and update their sent status.
- API to start/stop the auto sender.
- API to retrieve the list of sent messages.
- (Bonus) Save messageId to Redis after successful sending.

## Setup

### 1. Clone the repo
```bash
git clone <repo-url>
cd messenging_test
```

### 2. Configure environment variables
Create a `.env` file with the following variables:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=messaging
DB_SSLMODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
WEBHOOK_URL=https://webhook.site/your-webhook-id
WEBHOOK_KEY=your-webhook-key
```

### 3. Run with Docker
```bash
docker build -t messaging-app .
docker run --env-file .env -p 80:80 messaging-app
```

### 4. Run locally
```bash
go mod tidy
go run main.go
```

## API
- `POST /api/v1/auto-sender/start` — Start the auto sender
- `POST /api/v1/auto-sender/stop` — Stop the auto sender
- `GET /api/v1/messages/sent` — Get the list of sent messages
- Swagger UI: `/docs`

## Database
- Table `messages`: id, to, content, sent, sent_at, created_at, updated_at

## Bonus Redis
- After successful sending, save messageId and sent_at to Redis with key: `sent_message:{id}`

## Evaluation
- Meets architecture, clean code, documentation, and all requirements of the assignment.

# Introduction 
TODO: Give a short introduction of your project. Let this section explain the objectives or the motivation behind this project. 

# Getting Started
TODO: Guide users through getting your code up and running on their own system. In this section you can talk about:
1. Installation process
2. Software dependencies
3. Latest releases
4. API references

# Build and Test
TODO: Describe and show how to build your code and run the tests. 

# Contribute
TODO: Explain how other users and developers can contribute to make your code better. 

If you want to learn more about creating good readme files then refer the following [guidelines](https://docs.microsoft.com/en-us/azure/devops/repos/git/create-a-readme?view=azure-devops). You can also seek inspiration from the below readme files:
- [ASP.NET Core](https://github.com/aspnet/Home)
- [Visual Studio Code](https://github.com/Microsoft/vscode)
- [Chakra Core](https://github.com/Microsoft/ChakraCore)