# TestHeroBackendGo

A Go-based backend service for generating and managing standardized test questions using AI. This system supports multiple test types including ACT, SAT, MCAT, LSAT, GRE, and GMAT.

## Features

- AI-powered question generation
- Support for multiple test types and subjects
- LaTeX math expression support
- Markdown formatting for questions
- Adaptive difficulty based on user performance
- Docker containerization for easy deployment

## Prerequisites

- Go 1.19 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)
- OpenAI API key

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/TestHeroBackendGo.git
cd TestHeroBackendGo
```

2. Create a `.env` file in the root directory:
```env
OPENAI_API_KEY=your_api_key_here
DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=testhero
DB_PORT=5432
```

3. Install Go dependencies:
```bash
go mod download
```

## Development Setup

1. Start the PostgreSQL database and other services using Docker Compose:
```bash
docker-compose up --build
```
The server will start on `localhost:8080` by default.

## Project Structure

```
.
├── agent/              # AI agent and prompt management
├── controllers/        # HTTP request handlers
├── models/            # Database models and schemas
├── tasks/             # Topic data and processing
└── docker-compose.yml # Docker configuration
```

## API Endpoints

- `POST /api/questions/generate` - Generate a new question
- `POST /api/questions/similar/:questionId` - Generate a similar question
- `POST /api/questions/relevant` - Generate a question based on user performance

## Docker Deployment

The project includes a Docker Compose configuration for easy deployment:

```bash
# Build and start all services
docker-compose up --build

# Stop all services
docker-compose down
```

## Testing

Run the test suite:

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.