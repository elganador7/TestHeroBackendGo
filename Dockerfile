# Use an official Golang image as the base
FROM golang:1.22 as builder

# Set the working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Run the application
CMD ["go", "run", "main.go"]
