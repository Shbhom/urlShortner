import http from 'k6/http';
import { check } from 'k6';

export const options = {
    duration: '5m',
    vus: 200,
};

export default function () {
    const isHotKey = Math.random() < 0.90;
    let code;
    
    if (isHotKey) {
        code = 'seed_000000'; // The viral hot key
    } else {
        const randomId = Math.floor(Math.random() * 1000000);
        code = `seed_${randomId.toString().padStart(6, '0')}`;
    }
    
    const res = http.get(`http://localhost:8000/r/${code}`, { redirects: 0 });
    check(res, { 'status is 308': (r) => r.status === 308 });
}
