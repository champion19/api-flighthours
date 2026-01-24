import http from 'k6/http';
import { check, sleep, group } from 'k6';

// Test configuration
export const options = {
  stages: [
    { duration: '30s', target: 10 },  // Ramp-up to 10 users over 30s
    { duration: '1m', target: 10 },   // Stay at 10 users for 1 minute
    { duration: '30s', target: 20 },  // Ramp-up to 20 users
    { duration: '1m', target: 20 },   // Stay at 20 users for 1 minute
    { duration: '30s', target: 0 },   // Ramp-down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.05'],   // Error rate should be less than 5%
  },
};

const BASE_URL = 'http://host.docker.internal:8085/flighthours/api/v1';

export default function () {
  // Test 1: GET Messages List (Success)
  group('Messages API - Success Cases', function () {
    const messagesRes = http.get(`${BASE_URL}/messages`);

    check(messagesRes, {
      'GET /messages status is 200': (r) => r.status === 200,
      'GET /messages response time < 500ms': (r) => r.timings.duration < 500,
      'GET /messages has body': (r) => r.body && r.body.length > 0,
    });

    // Test 2: GET Message by ID (simulating existing ID)
    const messageByIdRes = http.get(`${BASE_URL}/messages/1`);

    check(messageByIdRes, {
      'GET /messages/:id response received': (r) => r.status === 200 || r.status === 404,
      'GET /messages/:id response time < 500ms': (r) => r.timings.duration < 500,
    });
  });

  // Test 3: Messages API - Error Cases
  group('Messages API - Error Cases', function () {
    // Test 404 Not Found
    const notFoundRes = http.get(`${BASE_URL}/messages/99999999`);

    check(notFoundRes, {
      'GET /messages/nonexistent status is 404': (r) => r.status === 404,
    });

    // Test malformed requests
    const badRequestRes = http.get(`${BASE_URL}/messages/invalid-id`);

    check(badRequestRes, {
      'GET /messages/invalid handles error': (r) => r.status >= 400,
    });
  });

  // Test 4: Filtering Messages
  group('Messages API - Filters', function () {
    const filters = [
      'module=users',
      'type=ERROR',
      'category=usuario_final',
      'active=true',
    ];

    filters.forEach((filter) => {
      const filterRes = http.get(`${BASE_URL}/messages?${filter}`);

      check(filterRes, {
        [`GET /messages?${filter} status is 200`]: (r) => r.status === 200,
        [`GET /messages?${filter} response time < 500ms`]: (r) => r.timings.duration < 500,
      });
    });
  });

  // Test 5: Prometheus Metrics Endpoint
  group('Prometheus Metrics', function () {
    const metricsRes = http.get('http://host.docker.internal:8085/metrics');

    check(metricsRes, {
      'GET /metrics status is 200': (r) => r.status === 200,
      'GET /metrics has prometheus format': (r) => r.body.includes('flighthours_http_requests_total'),
      'GET /metrics response time < 200ms': (r) => r.timings.duration < 200,
    });
  });

  // Simulate user think time
  sleep(1);
}

// Lifecycle hooks
export function setup() {
  console.log('Starting K6 load test for Flighthours Backend...');
  console.log(`Base URL: ${BASE_URL}`);
}

export function teardown(data) {
  console.log('K6 load test completed!');
}
