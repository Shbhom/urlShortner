import http from 'k6/http';
import { check } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    scenarios: {
        safe_load: {
            executor: 'shared-iterations',
            vus: 200,
            iterations: 1000000,
            maxDuration: '10m',
        },
    },
};

export default function () {
    const isInvalid = Math.random() < 0.90;
    let code;
    
    if (isInvalid) {
        code = `doesnotexist_${randomString(5)}`;
    } else {
        const randomId = Math.floor(Math.random() * 1000000);
        code = `seed_${randomId.toString().padStart(6, '0')}`;
    }
    
    const res = http.get(`http://localhost:8000/r/${code}`, { redirects: 0 });
    
    if (isInvalid) {
        check(res, { 'status is 404': (r) => r.status === 404 });
    } else {
        check(res, { 'status is 302': (r) => r.status === 302 });
    }
}
