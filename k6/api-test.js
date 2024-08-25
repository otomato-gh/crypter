// import necessary module
import http from "k6/http";
import { SharedArray } from 'k6/data';

export const options = {
  // define thresholds
  thresholds: {
    http_req_failed: ['rate<0.01'], // http errors should be less than 1%
    http_req_duration: ['p(99)<1000'], // 99% of requests should be below 1s
  },
  // define scenarios
  scenarios: {
    ramping_up: {
      executor: "ramping-vus",
      stages: [
        { duration: "5s", target: 20 },
        { duration: "5s", target: 60 },
        { duration: "10s", target: 100 },
      ],
    },
  }
};

const data = new SharedArray('texts', function () {
        const f = JSON.parse(open('./input.json'));
        console.log(typeof f);
        return f; // f must be an array[]
  });

export default function () {
  
  // define URL and payload
  var crypter_api_url = 'http://localhost:8080';
  if (__ENV.CRYPTER_API_URL === undefined) {
    console.log("CRYPTER_API_URL is not set, using default http://localhost:8080");

  } else {
    console.log(`CRYPTER_API_URL is set to ${__ENV.CRYPTER_API_URL}`);
    crypter_api_url = `${__ENV.CRYPTER_API_URL}`;
  }
  
  const url = crypter_api_url+'/encrypt';


  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  // send a post request and save response as a variable
  const randomText = data[Math.floor(Math.random() * data.length)];
  const payload = JSON.stringify(randomText);
  const res = http.post(url, payload, params);

  console.log(res.body);
  console.log(`CRYPTER_API_URL is set to ${__ENV.CRYPTER_API_URL}`);
}

