import http from 'k6/http';
import { check, sleep } from 'k6';

// Spike Test: Sudden traffic increase
export const options = {
  stages: [
    { duration: '1m', target: 10 },   // Normal load
    { duration: '10s', target: 100 }, // Sudden spike!
    { duration: '2m', target: 100 },  // Stay at spike
    { duration: '10s', target: 10 },  // Drop back to normal
    { duration: '1m', target: 10 },   // Normal load
    { duration: '10s', target: 0 },   // Ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // Higher threshold for spike
    http_req_failed: ['rate<0.15'],    // Allow more errors during spike
  },
};

const BASE_URL = 'http://host.docker.internal:8085/flighthours/api/v1';

export default function () {
  const messagesRes = http.get(`${BASE_URL}/messages`);

  check(messagesRes, {
    'status is 2xx': (r) => r.status >= 200 && r.status < 300,
    'response received': (r) => r.body !== null,
  });

  sleep(Math.random() * 2); // Random think time 0-2s
}

export function setup() {
  console.log('Running Spike Test - Simulating traffic spike...');
}

export function teardown(data) {
  console.log('Spike Test completed!');
}
