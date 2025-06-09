# SupplyChain Tracer API

A blockchain-powered supply chain tracking system built with Go, PostgreSQL, Redis, and Ethereum integration.

## 🚀 Features

- **Product Tracking**: Complete lifecycle tracking from manufacturer to end consumer
- **Blockchain Integration**: Immutable record storage on Ethereum blockchain
- **Multi-Stakeholder**: Support for manufacturers, distributors, retailers, and consumers
- **Real-time Updates**: WebSocket support for live tracking updates
- **Analytics Dashboard**: Supply chain metrics and compliance reporting
- **API Documentation**: Complete OpenAPI/Swagger documentation

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │       API       │    │   Blockchain    │
│   (Optional)    │◄──►│   (Go/Fiber)    │◄──►│   (Ethereum)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   PostgreSQL    │    │      Redis      │
                       │   (Main DB)     │    │    (Cache)      │
                       └─────────────────┘    └─────────────────┘
```

## 🛠️ Tech Stack

- **Backend**: Go 1.21+ with Gin framework
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Blockchain**: Ethereum (Sepolia testnet for development)
- **Queue**: Redis for async blockchain operations
- **Documentation**: Swagger/OpenAPI 3.0
- **Testing**: Testify, Ginkgo
- **Deployment**: Docker, Docker Compose

## 🚦 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)
- Ethereum wallet with testnet ETH

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/supplychain-tracer.git
   cd supplychain-tracer
   ```

2. **Setup environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Setup database**
   ```bash
   make migrate-up
   ```

5. **Run the application**
   ```bash
   make run
   ```

### Using Docker

```bash
# Development environment
docker-compose -f docker/docker-compose.dev.yml up

# Production environment
docker-compose -f docker/docker-compose.yml up
```

## 📖 API Documentation

API documentation is available at `http://localhost:8080/swagger/` when running the server.

### Main Endpoints

#### Products
- `POST /api/v1/products` - Create new product
- `GET /api/v1/products` - List products
- `GET /api/v1/products/{id}` - Get product details
- `PUT /api/v1/products/{id}` - Update product

#### Supply Chain Events
- `POST /api/v1/supply-chain/events` - Add tracking event
- `GET /api/v1/supply-chain/{productId}/history` - Get product history
- `GET /api/v1/supply-chain/events/{eventId}` - Get event details

#### Blockchain
- `POST /api/v1/blockchain/sync/{eventId}` - Sync event to blockchain
- `GET /api/v1/blockchain/verify/{hash}` - Verify blockchain transaction

#### Authentication
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/register` - Register stakeholder
- `GET /api/v1/auth/profile` - Get user profile

## 🧪 Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Run all tests with coverage
make test-coverage
```
## 🔧 Configuration

### Environment Variables

```env
# Server
PORT=8080
GIN_MODE=release

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=supplychain_tracer
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Blockchain
ETHEREUM_RPC_URL=https://sepolia.infura.io/v3/your-project-id
PRIVATE_KEY=your-private-key
CONTRACT_ADDRESS=0x...

# JWT
JWT_SECRET=your-jwt-secret
JWT_EXPIRE_HOURS=24

# External APIs
IPFS_URL=https://ipfs.infura.io:5001
```

## 🚀 Deployment

### Docker Deployment

```bash
# Build image
docker build -f docker/Dockerfile -t supplychain-tracer .

# Run with docker-compose
docker-compose -f docker/docker-compose.yml up -d
```

### Manual Deployment

```bash
# Build binary
make build

# Run migrations
make migrate-up

# Start server
./bin/supplychain-tracer
```

## 📈 Monitoring & Analytics

- Health check endpoint: `GET /health`
- Metrics endpoint: `GET /metrics` (Prometheus format)
- Analytics dashboard: `GET /api/v1/analytics/dashboard`

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Ethereum Foundation for blockchain infrastructure
- Go community for excellent libraries
- PostgreSQL team for robust database system

## 📞 Support

For support, email support@supplychain-tracer.com or join our Slack channel.

---

**Made with ❤️ by [Koriebruh]**