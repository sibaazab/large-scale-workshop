#!/bin/bash

# Function to handle SIGINT
cleanup() {
  echo "Caught SIGINT signal! Killing all started services..."
  
  kill 0
  exit 1
}

# Trap SIGINT signal
trap cleanup SIGINT


# Run the RunRegistry.sh script 3 times
for i in {1..3}; do
  echo "Starting RunRegistry.sh instance $i..."
  ./RunRegistry.sh &
  sleep 5
done

echo "Finished starting registry services..."

sleep 20

# Run the RunTestService.sh script 3 times
for i in {1..3}; do
  echo "Starting RunTestService.sh instance $i..."
  ./RunTestService.sh &
  sleep 20
done

echo "Finished starting test services..."
sleep 10

# Run the RunCache.sh script 3 times
for i in {1..3}; do
  echo "Starting RunCache.sh instance $i..."
  ./RunCache.sh &
  sleep 10
done

echo "Finished starting cache services..."
echo "All services have been started!"

wait