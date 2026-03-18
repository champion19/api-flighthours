import http from 'k6/http';
import { check, sleep } from 'k6';

export const BASE_URL = 'http://host.docker.internal:8085/flighthours/api/v1'; // NOSONAR — internal Docker test URL

export const defaultThresholds = {
  http_req_duration: ['p(95)<2000'],
  http_req_failed: ['rate<0.15'],
};

export function defaultScenario() {
  const messagesRes = http.get(`${BASE_URL}/messages`);

  check(messagesRes, {
    'status is 2xx': (r) => r.status >= 200 && r.status < 300,
    'response received': (r) => r.body !== null,
  });

  sleep(Math.random() * 2); // NOSONAR — test think-time simulation
}
