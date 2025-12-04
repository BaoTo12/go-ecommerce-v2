import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';
import { SharedArray } from 'k6/data';
import { crypto } from 'k6/experimental/webcrypto';

// Custom metrics
const purchaseAttempts = new Counter('flash_sale_attempts');
const purchaseSuccess = new Counter('flash_sale_success');
const purchaseFailed = new Counter('flash_sale_failed');
const soldOut = new Counter('flash_sale_sold_out');
const rateLimited = new Counter('flash_sale_rate_limited');
const powLatency = new Trend('pow_computation_time');
const successRate = new Rate('flash_sale_success_rate');

// Flash Sale Load Test - Simulating 11.11 scenario
export const options = {
    scenarios: {
        // Pre-sale warm-up
        warmup: {
            executor: 'constant-vus',
            vus: 10,
            duration: '30s',
            startTime: '0s',
            tags: { phase: 'warmup' },
        },
        // Spike at sale start (simulating midnight rush)
        spike: {
            executor: 'ramping-arrival-rate',
            startRate: 0,
            timeUnit: '1s',
            preAllocatedVUs: 1000,
            maxVUs: 10000,
            stages: [
                { duration: '10s', target: 5000 },  // Instant spike
                { duration: '60s', target: 10000 }, // Peak load
                { duration: '2m', target: 5000 },   // Sustained high
                { duration: '1m', target: 1000 },   // Wind down
            ],
            startTime: '1m',
            tags: { phase: 'spike' },
        },
    },
    thresholds: {
        http_req_duration: ['p(95)<1000'],        // 95% < 1s even under load
        flash_sale_success_rate: ['rate>0.01'],   // At least 1% success (many will fail due to stock)
        flash_sale_rate_limited: ['count<50000'], // Rate limiting should kick in
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const FLASH_SALE_ID = __ENV.FLASH_SALE_ID || 'flash-001';

// Proof of Work - solve hash puzzle (simulated)
async function solvePoW(challenge, difficulty) {
    const startTime = Date.now();
    let nonce = 0;
    const prefix = '0'.repeat(difficulty);

    while (true) {
        const data = challenge + nonce.toString();
        const hashBuffer = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(data));
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');

        if (hashHex.startsWith(prefix)) {
            powLatency.add(Date.now() - startTime);
            return nonce.toString();
        }

        nonce++;
        if (nonce > 1000000) {
            // Failsafe - in real scenario, PoW difficulty should be calibrated
            return nonce.toString();
        }
    }
}

export default async function () {
    const userId = 'user-' + __VU + '-' + __ITER;

    purchaseAttempts.add(1);

    // Step 1: Get PoW challenge
    let challengeRes = http.get(`${BASE_URL}/api/v1/flash-sale/${FLASH_SALE_ID}/challenge?user_id=${userId}`);

    if (!check(challengeRes, { 'got challenge': (r) => r.status === 200 })) {
        purchaseFailed.add(1);
        return;
    }

    let challenge;
    try {
        const body = JSON.parse(challengeRes.body);
        challenge = body.challenge;
    } catch {
        purchaseFailed.add(1);
        return;
    }

    // Step 2: Solve PoW (difficulty 4 = find hash starting with 0000)
    const nonce = await solvePoW(challenge, 4);

    // Step 3: Attempt purchase
    const purchasePayload = JSON.stringify({
        flash_sale_id: FLASH_SALE_ID,
        user_id: userId,
        quantity: 1,
        challenge: challenge,
        nonce: nonce,
    });

    let purchaseRes = http.post(`${BASE_URL}/api/v1/flash-sale/purchase`, purchasePayload, {
        headers: { 'Content-Type': 'application/json' },
    });

    // Analyze response
    if (purchaseRes.status === 200 || purchaseRes.status === 201) {
        purchaseSuccess.add(1);
        successRate.add(1);
    } else if (purchaseRes.status === 429) {
        // Rate limited
        rateLimited.add(1);
        successRate.add(0);
    } else if (purchaseRes.status === 410 || (purchaseRes.body && purchaseRes.body.includes('sold out'))) {
        // Sold out
        soldOut.add(1);
        successRate.add(0);
    } else {
        purchaseFailed.add(1);
        successRate.add(0);
    }

    // Small random delay to prevent thundering herd
    sleep(Math.random() * 0.5);
}

export function handleSummary(data) {
    const metrics = data.metrics;

    return {
        'flash_sale_test_summary.json': JSON.stringify(data, null, 2),
        stdout: `
=== Flash Sale Load Test Summary (11.11 Simulation) ===

Test Duration: ${data.state.testRunDurationMs}ms

Purchase Metrics:
  Total Attempts: ${metrics.flash_sale_attempts?.values?.count || 0}
  Successful: ${metrics.flash_sale_success?.values?.count || 0}
  Failed: ${metrics.flash_sale_failed?.values?.count || 0}
  Sold Out: ${metrics.flash_sale_sold_out?.values?.count || 0}
  Rate Limited: ${metrics.flash_sale_rate_limited?.values?.count || 0}

Response Times:
  Avg: ${(metrics.http_req_duration?.values?.avg || 0).toFixed(2)}ms
  P95: ${(metrics.http_req_duration?.values?.['p(95)'] || 0).toFixed(2)}ms
  P99: ${(metrics.http_req_duration?.values?.['p(99)'] || 0).toFixed(2)}ms

PoW Computation:
  Avg: ${(metrics.pow_computation_time?.values?.avg || 0).toFixed(2)}ms
  P95: ${(metrics.pow_computation_time?.values?.['p(95)'] || 0).toFixed(2)}ms

Success Rate: ${((metrics.flash_sale_success_rate?.values?.rate || 0) * 100).toFixed(2)}%
`,
    };
}
