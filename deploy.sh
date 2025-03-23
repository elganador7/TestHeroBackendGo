#!/bin/bash

# Exit on any error
set -e

# Configuration
PROJECT_ID=$(gcloud config get-value project)
REGION="us-central1"  # Change this if you're using a different region
IMAGE_NAME="testhero"
SERVICE_NAME="testhero"

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

# Check if required tools are installed
check_requirements() {
    print_step "Checking requirements..."
    
    if ! command -v gcloud &> /dev/null; then
        print_error "gcloud CLI is not installed. Please install it first."
        exit 1
    fi

    if ! command -v docker &> /dev/null; then
        print_error "docker is not installed. Please install it first."
        exit 1
    fi
}

# Authenticate with Google Cloud
authenticate() {
    print_step "Authenticating with Google Cloud..."
    
    # Check if already authenticated
    if ! gcloud auth print-access-token &> /dev/null; then
        gcloud auth login
    fi
    
    # Configure docker to use GCR
    gcloud auth configure-docker
}

# Build the Docker image
build_image() {
    print_step "Building Docker image..."
    
    IMAGE_TAG="gcr.io/${PROJECT_ID}/${IMAGE_NAME}:$(date +%Y%m%d-%H%M%S)"
    docker build -t "${IMAGE_TAG}" .
    
    print_success "Docker image built successfully: ${IMAGE_TAG}"
    echo "IMAGE_TAG=${IMAGE_TAG}" > .env.deploy
}

# Push the image to Google Container Registry
push_image() {
    print_step "Pushing image to Google Container Registry..."
    
    source .env.deploy
    docker push "${IMAGE_TAG}"
    
    print_success "Image pushed successfully"
}

# Deploy to Cloud Run
deploy_to_cloud_run() {
    print_step "Deploying to Cloud Run..."
    
    source .env.deploy
    
    gcloud run deploy ${SERVICE_NAME} \
        --image "${IMAGE_TAG}" \
        --platform managed \
        --region "${REGION}" \
        --allow-unauthenticated \
        --set-env-vars="ENVIRONMENT=production" \
        --add-cloudsql-instances="${PROJECT_ID}:${REGION}:testhero-db-instance" \
        --service-account=testhero-service@"${PROJECT_ID}".iam.gserviceaccount.com

    print_success "Deployed to Cloud Run successfully"
}

# Main deployment process
main() {
    print_step "Starting deployment process..."

    print_step "${PROJECT_ID}"
    
    check_requirements
    authenticate
    build_image
    push_image
    deploy_to_cloud_run
    
    print_success "Deployment completed successfully!"
    
    # Get the service URL
    SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} --region ${REGION} --format='value(status.url)')
    echo -e "\nService is available at: ${SERVICE_URL}"
}

# Run the deployment
main 