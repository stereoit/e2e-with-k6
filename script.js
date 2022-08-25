import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
    // http.get('https://localhost:8080');
    console.log("APP_HOST=", process.env.APP_HOST)
    http.get('cool-water-1296.fly.dev'); // make this env var
    sleep(1);
}

// export let options = {
//     ext: {
//         loadimpact: {
//             projectID: 3598518,
//             // Test runs with the same name groups test runs together
//             name: "End 2 end test from github actions"
//         }
//     }
// }