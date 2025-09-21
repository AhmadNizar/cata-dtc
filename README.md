# cata-dtc

Pokemon API with Go, MySQL, and Redis

## Setup Instructions

### Prerequisites
- Docker and Docker Compose
- VS Code with Dev Containers extension (recommended)

### Getting Started

1. **Environment Configuration**
   Copy the example environment file and configure your settings:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` file with your preferred configuration values.

2. **Development Container Setup**
   This project uses `.devcontainer` for a consistent development environment with MySQL and Redis.

   - Open project in VS Code
   - When prompted, click "Reopen in Container" or use Command Palette: `Remote-Containers: Reopen in Container`
   - Wait for the container to build and start (includes Go, MySQL, and Redis)

3. **Database Initialization**
   Once the dev container is running, execute these scripts in order:

   ```bash
   # Initialize MySQL database
   ./scripts/init_mysql.sh
   ```

   ```bash
   # Create database tables
   ./scripts/create_table.sh
   ```

   When prompted for database password, refer to the `DB_PASSWORD` value in your `.env` file.

4. **Part 3 SQL Answers**
   Database query answers for Part 3 can be found in:
   - `/scripts/part3_1st_answer.sql`
   - `/scripts/part3_2nd_andwer.sql`

### Running the Application

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### Services
- **API**: Port 8080
- **MySQL**: Port 3306
- **Redis**: Port 6379

All services are automatically configured and accessible within the dev container.