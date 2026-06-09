import http from 'k6/http';
import { check } from 'k6';

export const options = {
    stages: [
        { duration: '1m', target: 50 },
        { duration: '1m', target: 100 },
        { duration: '1m', target: 250 },
        { duration: '1m', target: 500 },
        { duration: '1m', target: 1000 },
        { duration: '1m', target: 2000 },
        { duration: '1m', target: 5000 },
    ],
    thresholds: {
        // Test auto-aborts if p(95) is over 100ms or error rate exceeds 1%
        http_req_duration: ['p(95)<100'],
        http_req_failed: ['rate<0.01'], 
    },
};

export default function () {
    const randomId = Math.floor(Math.random() * 1000000);
    const code = `seed_${randomId.toString().padStart(6, '0')}`;
    
    const res = http.get(`http://localhost:8000/r/${code}`, { redirects: 0 });
    check(res, { 'status is 308': (r) => r.status === 308 });
}
