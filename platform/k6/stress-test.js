import { defaultScenario, defaultThresholds } from './shared-config.js';

export const options = {
  stages: [
    { duration: '2m', target: 50 },
    { duration: '5m', target: 50 },
    { duration: '2m', target: 100 },
    { duration: '5m', target: 100 },
    { duration: '2m', target: 0 },
  ],
  thresholds: defaultThresholds,
};

export default function () {
  defaultScenario();
}

export function setup() {
  console.log('Running Stress Test - Sustained high load...');
}

export function teardown(data) {
  console.log('Stress Test completed!');
}
