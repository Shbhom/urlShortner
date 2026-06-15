import http from 'k6/http';
import { check } from 'k6';

export const options = {
    scenarios: {
        safe_load: {
            executor: 'shared-iterations',
            vus: 200,
            iterations: 300000,
            maxDuration: '10m',
        },
    },
    thresholds: {
        http_req_duration: ['p(95)<100'],
        http_req_failed: ['rate<0.01'],
    },
};

export default function () {
    const randomId = Math.floor(Math.random() * 1000000);
    const code = `seed_${randomId.toString().padStart(6, '0')}`;
    
    const res = http.get(`http://localhost:8000/r/${code}`, { redirects: 0 });
    check(res, { 'status is 302': (r) => r.status === 302 });
}
