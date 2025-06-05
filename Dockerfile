# Use an official Golang image as the base
FROM golang:1.23 as builder

# Set the working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["go", "run", "main.go"]
