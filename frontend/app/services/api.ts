const API_BASE_URL = 'http://localhost:8080/api/v1';// dev

export const createPoll = async (title: string, options: string[]) => {
  const response = await fetch(`${API_BASE_URL}/polls`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ title, options }),
  });
  if (!response.ok) {
    throw new Error('Failed to create poll');
  }
  return response.json();
};

export const getPoll = async (id: string) => {
  const response = await fetch(`${API_BASE_URL}/polls/${id}/results`);
   if (!response.ok) {
    throw new Error('Poll not found');
  }
  return response.json();
};

export const voteOnPoll = async (id: string, optionIndex: number) => {
  const response = await fetch(`${API_BASE_URL}/polls/${id}/vote`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ option_indices: [optionIndex] }),
  });
  if (!response.ok) {
    throw new Error('Failed to vote');
  }
  return response.json();
};