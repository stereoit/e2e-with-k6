import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
    http.get('https://localhost:8080');
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