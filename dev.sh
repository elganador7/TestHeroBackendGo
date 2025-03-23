#!/bin/bash

# Exit on any error
set -e

# Configuration
PROJECT_ID=$(gcloud config get-value project)
REGION="us-central1"
INSTANCE_NAME="testhero-db-instance"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print step message
print_step() {
    echo -e "${YELLOW}[STEP]${NC} $1"
}

# Print success message
print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Print error message
print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Cloud SQL proxy is installed
check_proxy() {
    if ! command -v cloud-sql-proxy &> /dev/null; then
        print_error "Cloud SQL proxy is not installed. Installing now..."
        if [[ "$OSTYPE" == "darwin"* ]]; then
            brew install cloud-sql-proxy
        else
            # For Linux and other systems
            print_error "Please install cloud-sql-proxy manually from: https://cloud.google.com/sql/docs/postgres/connect-admin-proxy"
            exit 1
        fi
    fi
}

# Start Cloud SQL proxy
start_proxy() {
    print_step "Starting Cloud SQL proxy..."
    
    INSTANCE_CONNECTION_NAME="${PROJECT_ID}:${REGION}:${INSTANCE_NAME}"
    
    # Kill any existing cloud_sql_proxy process
    pkill -f cloud-sql-proxy || true
    
    # Start the proxy in the background
    cloud-sql-proxy "${INSTANCE_CONNECTION_NAME}" &
    
    # Wait for proxy to start
    sleep 5
    
    print_success "Cloud SQL proxy started successfully"
}

# Set up environment variables
setup_env() {
    print_step "Setting up environment variables..."
    
    # Create .env file if it doesn't exist
    if [ ! -f .env ]; then
        cat > .env << EOF
DB_HOST=localhost
DB_USER=testhero_app
DB_NAME=testhero
DB_PORT=5432
ENVIRONMENT=development
EOF
        print_step "Created .env file. Please edit it to add your database password and OpenAI API key"
    fi
    
    print_success "Environment setup complete"
}

# Start the application
start_app() {
    print_step "Starting the application..."
    
    go run main.go
}

# Main development process
main() {
    print_step "Starting development environment..."
    
    check_proxy
    setup_env
    start_proxy
    start_app
}

# Cleanup function
cleanup() {
    print_step "Cleaning up..."
    pkill -f cloud-sql-proxy || true
    print_success "Cleanup complete"
}

# Set up cleanup trap
trap cleanup EXIT

# Run the development environment
main 