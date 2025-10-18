#!/bin/bash

echo "=== FX Data Analysis Application ==="
echo ""
echo "Starting backend server..."
echo ""

# Start backend in background
go run cmd/server/main.go &
BACKEND_PID=$!

# Wait for backend to start
sleep 2

echo ""
echo "Backend server started (PID: $BACKEND_PID)"
echo "Backend URL: http://localhost:8080"
echo ""
echo "To start frontend, open a new terminal and run:"
echo "  cd frontend"
echo "  npm run dev"
echo ""
echo "Frontend URL: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop the backend server"
echo ""

# Wait for Ctrl+C
trap "kill $BACKEND_PID; echo 'Backend server stopped'; exit 0" INT TERM
wait $BACKEND_PID
