import http from 'k6/http';
import { check } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    duration: '5m',
    vus: 200,
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
        check(res, { 'status is 308': (r) => r.status === 308 });
    }
}
