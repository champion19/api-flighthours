import { defaultScenario, defaultThresholds } from './shared-config.js';

export const options = {
  stages: [
    { duration: '1m', target: 10 },
    { duration: '10s', target: 100 },
    { duration: '2m', target: 100 },
    { duration: '10s', target: 10 },
    { duration: '1m', target: 10 },
    { duration: '10s', target: 0 },
  ],
  thresholds: defaultThresholds,
};

export default function spikeTest() {
  defaultScenario();
}

export function setup() {
  console.log('Running Spike Test - Simulating traffic spike...');
}

export function teardown(data) {
  console.log('Spike Test completed!');
}
