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
    "title": "What programming languages do you know?",
    "options": ["Go", "Python", "JavaScript", "Rust", "TypeScript"]
  }')

echo "$RESPONSE" | jq .

# Extract poll_id
POLL_ID=$(echo "$RESPONSE" | jq -r .poll_id)
echo -e "\nCreated poll with ID: $POLL_ID\n"

# 3. Single choice voting
echo "3️⃣  Single choice voting"

echo "Vote for option 0 (Go)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_indices": [0]}' | jq .

echo "Vote for option 1 (Python)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_indices": [1]}' | jq .

echo "Vote for option 2 (JavaScript)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_indices": [2]}' | jq .

echo ""

# 4. Multiple choice voting
echo "4️⃣  Multiple choice voting"

echo "Vote for Go and Python (0, 1)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_indices": [0, 1]}' | jq .

echo "Vote for JavaScript, Rust, and TypeScript (2, 3, 4)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_indices": [2, 3, 4]}' | jq .

echo "Vote for all languages (0, 1, 2, 3, 4)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_indices": [0, 1, 2, 3, 4]}' | jq .

echo ""

# 5. Random voting
echo "5️⃣  Random voting (mixed single and multiple choice)"

for i in {1..5}; do
  # Generate random number of choices (1-3)
  NUM_CHOICES=$((RANDOM % 3 + 1))

  # Generate random unique indices
  INDICES=()
  while [ ${#INDICES[@]} -lt $NUM_CHOICES ]; do
    NEW_INDEX=$((RANDOM % 5))
    # Check if index already exists
    if [[ ! " ${INDICES[@]} " =~ " ${NEW_INDEX} " ]]; then
      INDICES+=($NEW_INDEX)
    fi
  done

  # Sort indices
  IFS=$'\n' SORTED_INDICES=($(sort -n <<<"${INDICES[*]}"))
  unset IFS

  # Convert to JSON array format
  JSON_INDICES=$(printf '%s\n' "${SORTED_INDICES[@]}" | jq -s .)

  echo "Random vote #$i: options $JSON_INDICES"
  curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
    -H "Content-Type: application/json" \
    -d "{\"option_indices\": $JSON_INDICES}" | jq -c .
done

echo ""

# 6. Checking results
echo "6️⃣  Poll results"
curl -s "${BASE_URL}/api/v1/polls/${POLL_ID}/results" | jq .
echo ""

# 7. Error cases
echo "7️⃣  Error cases"

echo "Trying to vote with empty array (should fail)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_indices": []}' | jq .
echo ""

echo "Trying to vote with duplicate indices (should fail)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_indices": [0, 0, 1]}' | jq .
echo ""

echo "Trying to vote for non-existent option 999 (should fail)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_indices": [999]}' | jq .
echo ""

echo "Trying to vote for negative option -1 (should fail)"
curl -s -X POST "${BASE_URL}/api/v1/polls/${POLL_ID}/vote" \
  -H "Content-Type: application/json" \
  -d '{"option_indices": [-1]}' | jq .
echo ""

# 8. Try to get non-existent poll
echo "8️⃣  Try to get non-existent poll (should fail)"
curl -s "${BASE_URL}/api/v1/polls/non-existent-id/results" | jq .
echo ""

echo "✅ All example requests are executed!"
echo ""
echo "🔗 Useful links:"
echo "- Vote: ${BASE_URL}/api/v1/polls/${POLL_ID}/vote"
echo "- Results: ${BASE_URL}/api/v1/polls/${POLL_ID}/results"
echo "- Swagger UI: ${BASE_URL}/swagger/index.html"
echo ""
echo "📝 Example requests:"
echo ""
echo "# Single choice:"
echo "curl -X POST ${BASE_URL}/api/v1/polls/${POLL_ID}/vote \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"option_indices\": [0]}'"
echo ""
echo "# Multiple choice:"
echo "curl -X POST ${BASE_URL}/api/v1/polls/${POLL_ID}/vote \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"option_indices\": [0, 2, 4]}'"