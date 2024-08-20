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
    // arbitrary name of scenario
    average_load: {
      executor: "ramping-vus",
      stages: [
        // ramp up to average load of 20 virtual users
        { duration: "10s", target: 100 },
        // maintain load
        { duration: "50s", target: 300 },
        // ramp down to zero
        { duration: "5s", target: 0 },
      ],
    },
    breaking: {
      executor: "ramping-vus",
      stages: [
        { duration: "10s", target: 20 },
        { duration: "50s", target: 20 },
        { duration: "50s", target: 40 },
        { duration: "50s", target: 60 },
        { duration: "50s", target: 80 },
        { duration: "50s", target: 100 },
        { duration: "50s", target: 120 },
        { duration: "50s", target: 140 },
        //....
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
  const url = "http://gocrypt:8080/encrypt";
//  const payload = JSON.stringify({
//    plaintext: "In this tutorial, you've used k6 to make a POST request and check that it responds with a 200 status.",
//  });


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
}

