#!/bin/bash

# WorkKG Project Deployment Script
# Usage: ./deploy.sh [both|front|back]

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -d "work_kg_backend" ] || [ ! -d "work_kg_frontend" ]; then
    print_error "This script must be run from the project root directory containing work_kg_backend and work_kg_frontend folders"
    exit 1
fi

# Function to force update Git repository (overwrite local changes)
force_git_update() {
    branch=$(git branch --show-current || echo "master")
    print_status "Fetching latest changes and overwriting local changes..."
    git fetch --all
    git reset --hard origin/$branch || git reset --hard origin/master
    git clean -fd
    print_success "Local changes cleared and updated to remote."
}

# Function to deploy backend (Go)
deploy_backend() {
    echo "🔧 Starting Backend Deployment..."
    echo "=================================="

    # Stop and delete existing PM2 process
    print_status "Stopping existing backend PM2 process..."
    pm2 stop work_kg_backend 2>/dev/null || print_warning "No existing backend process to stop"
    pm2 delete work_kg_backend 2>/dev/null || print_warning "No existing backend process to delete"

    # Navigate to backend directory
    cd work_kg_backend

    # Install Go dependencies
    print_status "Downloading Go dependencies..."
    go mod tidy

    # Build the Go binary
    print_status "Building Go binary..."
    go build -o server .

    # Start with PM2
    print_status "Starting backend with PM2..."
    pm2 start ./server --name "work_kg_backend"

    # Wait a moment for startup
    sleep 3

    # Check if it's running
    if pm2 list | grep -q "work_kg_backend.*online"; then
        print_success "Backend deployed successfully on port 7041!"
        print_status "Testing backend health..."
        sleep 2
        curl -s http://localhost:7041/api/stats && echo "" || print_warning "Health check failed"
    else
        print_error "Backend deployment failed!"
        print_status "Checking logs..."
        pm2 logs work_kg_backend --lines 10 --nostream
        exit 1
    fi

    # Go back to project root
    cd ..

    print_success "Backend deployment completed!"
    echo "=================================="
}

# Function to deploy frontend (Next.js)
deploy_frontend() {
    echo "🎨 Starting Frontend Deployment..."
    echo "=================================="

    # Stop and delete existing PM2 process
    print_status "Stopping existing frontend PM2 process..."
    pm2 stop work_kg_frontend 2>/dev/null || print_warning "No existing frontend process to stop"
    pm2 delete work_kg_frontend 2>/dev/null || print_warning "No existing frontend process to delete"

    # Navigate to frontend directory
    cd work_kg_frontend

    # Install dependencies
    print_status "Installing dependencies with npm..."
    npm install

    # Build the project
    print_status "Building Next.js project..."
    npm run build

    # Start with PM2
    print_status "Starting frontend with PM2..."
    PORT=7040 pm2 start npm --name "work_kg_frontend" -- start

    # Wait a moment for startup
    sleep 5

    # Check if it's running
    if pm2 list | grep -q "work_kg_frontend.*online"; then
        print_success "Frontend deployed successfully on port 7040!"
        print_status "Testing frontend..."
        sleep 2
        curl -s -I http://localhost:7040 | head -1 && echo "" || print_warning "Frontend check failed"
    else
        print_error "Frontend deployment failed!"
        print_status "Checking logs..."
        pm2 logs work_kg_frontend --lines 10 --nostream
        exit 1
    fi

    # Go back to project root
    cd ..

    print_success "Frontend deployment completed!"
    echo "=================================="
}

# Function to deploy both
deploy_both() {
    echo "🚀 Starting Full Deployment (Both Frontend & Backend)..."
    echo "======================================================="

    # Force update from git first
    force_git_update

    # Deploy backend first
    deploy_backend

    echo ""

    # Deploy frontend
    deploy_frontend

    # Save PM2 processes
    pm2 save

    # Final status check
    echo ""
    print_status "Checking final deployment status..."
    sleep 3

    pm2 status | grep work_kg

    print_success "🎉 Full deployment completed!"
    print_success "Frontend: http://localhost:7040 (work_kg.okugula.dev)"
    print_success "Backend: http://localhost:7041 (work_kg_backend.okugula.dev)"
    echo "======================================================="
}

# Main script logic
case "${1:-both}" in
    "both")
        deploy_both
        ;;
    "front")
        deploy_frontend
        ;;
    "back")
        deploy_backend
        ;;
    *)
        echo "🚀 WorkKG Project Deployment Script"
        echo "===================================="
        echo "Usage: ./deploy.sh [both|front|back]"
        echo ""
        echo "Options:"
        echo "  both   - Deploy both frontend and backend (default)"
        echo "  front  - Deploy frontend only"
        echo "  back   - Deploy backend only"
        echo ""
        echo "Examples:"
        echo "  ./deploy.sh        # Deploy both"
        echo "  ./deploy.sh both   # Deploy both"
        echo "  ./deploy.sh front  # Deploy frontend only"
        echo "  ./deploy.sh back   # Deploy backend only"
        exit 1
        ;;
esac
