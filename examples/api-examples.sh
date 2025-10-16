#!/bin/bash

# API Examples for Poll Service
# Check that service starts at http://localhost:8080

set -e

BASE_URL="http://localhost:8080"

echo "🚀 Poll Service API Examples"
echo "=============================="
echo ""

# 1. Health Check
echo "1️⃣  Health Check"
curl -s "${BASE_URL}/health" | jq .
echo -e "\n"

# 2. Create a poll
echo "2️⃣  Creating poll"
RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/polls" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Favorite programming language",
    "options": ["Go", "Python", "JavaScript", "Rust", "TypeScript"]
  }')

echo "$RESPONSE" | jq .

# Extract poll_id
POLL_ID=$(echo "$RESPONSE" | jq -r .poll_id)
echo -e "\nCreated poll with ID: $POLL_ID\n"

# 3. Voting
echo "3️⃣  Voting"

for i in {0..4}; do
  OPTION=$((RANDOM % 5))
  echo "Vote for $OPTION"
  curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
    -H "Content-Type: application/json" \
    -d "{\"option_index\": $OPTION}" | jq .
done
echo ""

echo "Add more votes..."
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_index": 0}' > /dev/null

curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_index": 0}' > /dev/null

curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_index": 1}' > /dev/null

echo ""

# 4. Checking results
echo "4️⃣  Poll results"
curl -s "${BASE_URL}/api/v1/polls/${POLL_ID}/results" | jq .
echo ""

# 5. Vote for non-valid option
echo "5️⃣  Trying vote for non-exist option (expected error)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_index": 999}' | jq .
echo ""

# 6. Try to get non-exist poll
echo "6️⃣  Try to get non-exist poll (expected error)"
curl -s "${BASE_URL}/api/v1/polls/non-existent-id/results" | jq .
echo ""

echo "✅ All example requests are executed!"
echo ""
echo "Useful links:"
echo "- Poll: ${BASE_URL}/api/v1/polls/${POLL_ID}/vote"
echo "- Results: ${BASE_URL}/api/v1/polls/${POLL_ID}/results"
