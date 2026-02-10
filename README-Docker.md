# Docker Configuration Guide

This project provides optimized Docker configurations for development, staging, and production environments.

## Quick Start

### Development Environment

```bash
# Start all services including development tools
docker-compose --profile dev up -d

# Start only the API and external services (without dev tools)
docker-compose up -d
```

### Staging Environment

```bash
# Deploy to staging with external services
docker-compose -f docker-compose.staging.yml --env-file .env.staging up -d
```

## Environment Configurations

### Development

- **File**: `docker-compose.yml`
- **Services**: API, PostgreSQL, Redis, MinIO, PgAdmin, Redis Commander
- **Network**: `app-network` + `traefik_proxy`
- **Profiles**: `dev` (optional development tools)

### Staging

- **File**: `docker-compose.staging.yml`
- **Services**: API only (connects to external services)
- **Network**: `traefik_proxy` only
- **Health Checks**: PostgreSQL, Redis, MinIO

## Traefik Integration

The API service is configured to work with Traefik reverse proxy:

### Production

- **Domain**: `api.yourdomain.com`
- **Port**: 3098 (internal)
- **TLS**: Let's Encrypt automatic certificates

### Staging

- **Domain**: `api-staging.yourdomain.com`
- **Port**: 3098 (internal)
- **TLS**: Let's Encrypt automatic certificates

## External Services Configuration

When deploying to staging/production with external services:

### Environment Variables Required

```bash
# Database
DATABASE_URL=postgres://user:pass@host:port/db?sslmode=disable

# Redis
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password

# MinIO
MINIO_ENDPOINT=your-minio-host:9000
MINIO_ACCESS_KEY=your-access-key
MINIO_SECRET_KEY=your-secret-key
MINIO_BUCKET=your-bucket-name
```

### Health Check Services

The staging configuration includes health check services that ensure external dependencies are ready before starting the API:

- **PostgreSQL Health Check**: Verifies database connectivity
- **Redis Health Check**: Verifies Redis connectivity
- **MinIO Health Check**: Verifies object storage connectivity

## Security Features

### Dockerfile Security

- Non-root user (`appuser`)
- Minimal Alpine Linux base
- Optimized binary with stripped symbols
- Proper signal handling with ENTRYPOINT

### Network Security

- Dedicated Docker networks
- External service isolation
- Traefik integration for SSL termination

## Volume Management

### Development

- `pgdata`: PostgreSQL data persistence
- `redis-data`: Redis data persistence
- `minio-data`: MinIO object storage
- `pgadmin-data`: PgAdmin configuration

### Production/Staging

- No persistent volumes (data managed by external services)
- Logs volume for application logging

## Monitoring and Debugging

### Development Tools

- **PgAdmin**: Database management at `http://localhost:5050`
- **Redis Commander**: Redis management at `http://localhost:8081`

### Health Checks

All services include health checks for monitoring:

- PostgreSQL: `pg_isready` command
- Redis: `redis-cli ping` command
- MinIO: HTTP health endpoint

## Deployment Commands

### Development

```bash
# Full development stack
docker-compose --profile dev up -d

# Stop development stack
docker-compose down

# View logs
docker-compose logs -f api
```

### Staging

```bash
# Deploy staging
docker-compose -f docker-compose.staging.yml --env-file .env.staging up -d

# Stop staging
docker-compose -f docker-compose.staging.yml down

# View staging logs
docker-compose -f docker-compose.staging.yml logs -f api
```

### Production

```bash
# Deploy production (replace with your production compose file)
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d
```

## Troubleshooting

### Common Issues

1. **Port Conflicts**: Ensure ports 3098, 5432, 6379, 9000, 9001 are available
2. **Network Issues**: Ensure `traefik_proxy` network exists
3. **Permission Issues**: Check file permissions for mounted volumes
4. **Health Check Failures**: Verify external service connectivity

### Debug Commands

```bash
# Check service status
docker-compose ps

# Check health status
docker-compose ps --health

# View specific service logs
docker-compose logs api
docker-compose logs postgres

# Execute commands in container
docker-compose exec api sh
docker-compose exec postgres psql -U postgres
```

## Performance Optimizations

### Build Optimizations

- Layer caching with separate go.mod copy
- Optimized Go build flags (`-ldflags='-w -s'`)
- Multi-stage build for minimal runtime image

### Runtime Optimizations

- Alpine Linux for smaller footprint
- Non-root user for security
- Proper resource limits (add to production compose)

### Network Optimizations

- Dedicated networks for service isolation
- Traefik for efficient load balancing
- Health checks for service reliability
