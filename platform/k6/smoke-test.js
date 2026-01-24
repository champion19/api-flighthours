import http from 'k6/http';
import { check, sleep } from 'k6';

// Smoke Test: Minimal load to verify system works
export const options = {
  vus: 1,
  duration: '30s',
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'], // Less than 1% errors
  },
};

const BASE_URL = 'http://host.docker.internal:8085/flighthours/api/v1';

export default function () {
  // Basic health check
  const messagesRes = http.get(`${BASE_URL}/messages`);

  check(messagesRes, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  // Check Prometheus metrics
  const metricsRes = http.get('http://host.docker.internal:8085/metrics');

  check(metricsRes, {
    'metrics endpoint works': (r) => r.status === 200,
  });

  sleep(1);
}

export function setup() {
  console.log('Running Smoke Test...');
}

export function teardown(data) {
  console.log('Smoke Test completed!');
}
