import http from 'k6/http';
import { check } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    scenarios: {
        safe_load: {
            executor: 'shared-iterations',
            vus: 50,
            iterations: 100000,
            maxDuration: '10m',
        },
    },
};

export default function () {
    const payload = JSON.stringify({
        url: `https://example.com/dynamic-${randomString(10)}`
    });
    
    const params = {
        headers: { 'Content-Type': 'application/json' },
    };
    
    const res = http.post('http://localhost:8000/url', payload, params);
    check(res, { 'status is 201': (r) => r.status === 201 });
}
