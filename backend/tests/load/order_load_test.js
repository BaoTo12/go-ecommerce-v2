import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

// Custom metrics
const orderCreated = new Counter('orders_created');
const orderFailed = new Counter('orders_failed');
const orderLatency = new Trend('order_latency');
const successRate = new Rate('success_rate');

// Test configuration
export const options = {
  scenarios: {
    // Smoke test - verify basic functionality
    smoke: {
      executor: 'constant-vus',
      vus: 1,
      duration: '30s',
      startTime: '0s',
      tags: { test_type: 'smoke' },
    },
    // Load test - normal load
    load: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 100 },  // Ramp up
        { duration: '5m', target: 100 },  // Stay at 100
        { duration: '2m', target: 0 },    // Ramp down
      ],
      startTime: '1m',
      tags: { test_type: 'load' },
    },
    // Stress test - find breaking point
    stress: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 100 },
        { duration: '5m', target: 200 },
        { duration: '5m', target: 500 },
        { duration: '5m', target: 1000 },
        { duration: '5m', target: 500 },
        { duration: '2m', target: 0 },
      ],
      startTime: '10m',
      tags: { test_type: 'stress' },
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // 95% < 500ms, 99% < 1s
    success_rate: ['rate>0.95'],                     // 95% success rate
    orders_failed: ['count<100'],                    // Max 100 failed orders
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Test data
const products = [
  { id: 'prod-001', name: 'Laptop', price: 999.99 },
  { id: 'prod-002', name: 'Phone', price: 699.99 },
  { id: 'prod-003', name: 'Headphones', price: 199.99 },
  { id: 'prod-004', name: 'Watch', price: 299.99 },
  { id: 'prod-005', name: 'Tablet', price: 499.99 },
];

function getRandomProduct() {
  return products[Math.floor(Math.random() * products.length)];
}

function generateUserId() {
  return 'user-' + Math.random().toString(36).substring(7);
}

export default function () {
  const userId = generateUserId();
  
  group('Order Creation Flow', () => {
    // Step 1: Add to cart
    const product = getRandomProduct();
    const cartPayload = JSON.stringify({
      user_id: userId,
      product_id: product.id,
      product_name: product.name,
      quantity: Math.floor(Math.random() * 3) + 1,
      unit_price: product.price,
    });

    let cartRes = http.post(`${BASE_URL}/api/v1/cart/add`, cartPayload, {
      headers: { 'Content-Type': 'application/json' },
    });

    check(cartRes, {
      'cart add status 200': (r) => r.status === 200,
    });

    // Step 2: Get cart
    let getCartRes = http.get(`${BASE_URL}/api/v1/cart/${userId}`);
    check(getCartRes, {
      'get cart status 200': (r) => r.status === 200,
    });

    // Step 3: Create order
    const startTime = Date.now();
    
    const orderPayload = JSON.stringify({
      user_id: userId,
      shipping_address: '123 Test Street, City, Country',
      payment_method_id: 'pm_test_' + Math.random().toString(36).substring(7),
    });

    let orderRes = http.post(`${BASE_URL}/api/v1/checkout/initiate`, orderPayload, {
      headers: { 'Content-Type': 'application/json' },
    });

    const endTime = Date.now();
    orderLatency.add(endTime - startTime);

    const orderSuccess = check(orderRes, {
      'order creation status 200': (r) => r.status === 200 || r.status === 201,
      'order has session_id': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.session_id !== undefined;
        } catch {
          return false;
        }
      },
    });

    if (orderSuccess) {
      orderCreated.add(1);
      successRate.add(1);
    } else {
      orderFailed.add(1);
      successRate.add(0);
    }
  });

  sleep(Math.random() * 2 + 1); // 1-3 seconds between iterations
}

export function handleSummary(data) {
  return {
    'order_load_test_summary.json': JSON.stringify(data, null, 2),
    stdout: textSummary(data, { indent: ' ', enableColors: true }),
  };
}

function textSummary(data, options) {
  const metrics = data.metrics;
  return `
=== K6 Load Test Summary ===

Duration: ${data.state.testRunDurationMs}ms
VUs Max: ${data.root_group.checks.length}

HTTP Requests:
  Total: ${metrics.http_reqs?.values?.count || 0}
  Rate: ${(metrics.http_reqs?.values?.rate || 0).toFixed(2)}/s

Response Times:
  Avg: ${(metrics.http_req_duration?.values?.avg || 0).toFixed(2)}ms
  P95: ${(metrics.http_req_duration?.values?.['p(95)'] || 0).toFixed(2)}ms
  P99: ${(metrics.http_req_duration?.values?.['p(99)'] || 0).toFixed(2)}ms

Orders:
  Created: ${metrics.orders_created?.values?.count || 0}
  Failed: ${metrics.orders_failed?.values?.count || 0}
  Success Rate: ${((metrics.success_rate?.values?.rate || 0) * 100).toFixed(2)}%
`;
}
