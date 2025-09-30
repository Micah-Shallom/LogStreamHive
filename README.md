# LogStreamHive

<img width="1382" height="672" alt="image" src="https://github.com/user-attachments/assets/56bd127a-f2cb-4550-98d7-ee728d92a5bb" />

## Features

- 🚀 Real-time log streaming and visualization
- 🔍 Advanced log filtering and search capabilities
- 📊 Interactive dashboard with metrics
- 🔐 Configurable logging levels and formats
- 🐳 Docker-ready deployment
- ⚡ Built with Next.js and Go for high performance

## Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- Docker (optional)

### Running Locally

1. Start the logger service:
```bash
cd logger
go run main.go
```

2. Launch the web interface:
```bash
cd client
npm install
npm run dev
```

3. Access the web interface at `http://localhost:3000`

### Docker Deployment

```bash
docker compose up -d
```

## Configuration

### Logger Service

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `PORT` | Service port | `8080` |
| `OUTPUT_FORMAT` | Log output format (json, text) | `json` |

### Web Interface

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Logger service URL | `http://localhost:8080` |
| `PORT` | Web interface port | `3000` |

## API Documentation

- `GET /api/logs`: Fetch logs with optional filters
- `POST /api/logs`: Submit new log entries
- `GET /api/stats`: Retrieve logging statistics

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
