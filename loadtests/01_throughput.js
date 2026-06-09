import http from 'k6/http';
import { check } from 'k6';

export const options = {
    stages: [
        { duration: '2m', target: 50 },
        { duration: '2m', target: 100 },
        { duration: '2m', target: 250 },
        { duration: '2m', target: 500 },
        { duration: '2m', target: 1000 },
    ],
};

export default function () {
    const randomId = Math.floor(Math.random() * 1000000);
    const code = `seed_${randomId.toString().padStart(6, '0')}`;
    
    const res = http.get(`http://localhost:8000/r/${code}`, { redirects: 0 });
    check(res, { 'status is 308': (r) => r.status === 308 });
}
