import http from 'k6/http';

export const options = {
  vus: 1,
  duration: '20s',
};

export default function () {
  http.post(
    'http://localhost:8080/wallet',
    JSON.stringify({
    "wallet_id": 1,
    "type": "deposit",
    "amount": 1   }),
    { headers: { 'Content-Type': 'application/json' } }
  );
}