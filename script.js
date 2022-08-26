import http from 'k6/http';
import { sleep, group, check } from 'k6';

const options = {
    vus: 10,
    duration: '60s',
};
const SLEEP_DURATION = 0.1;

const baseUrl = __ENV.APP_HOST;

export default function () {
    const res = http.get(`http://${baseUrl}/`);
    sleep(1);

    const params = {
        headers: {
            'Content-Type': 'application/json',
        }
    };

    // const login_response = login();
    // check(login_response, {
    //     'is status 200': (r) => r.status === 200,
    //     'is access token present': (r) => r.json().hasOwnProperty('accessToken'),
    // })

    // params['Auth'] = `Bearer ${login_response.json()['accessToken']}`;
    sleep(SLEEP_DURATION);

    group('Articles', (_) => {
        params.tags.name = 'get-all-articles'
        // Get all articles
        const articles_response = http.get(
            `${baseUrl}/articles`, params
        );
        check(articles_response, {
            'is status 200': (r) => r.status === 200,
        })
        sleep(SLEEP_DURATION);


        // Create new article
        // Get article details
        // Delete article
    });
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


function login() {
    const url = 'http://api.platform.com/login';

    const payload = JSON.stringify({
        email: 'robert.smol@2k.com',
        password: 'PASSWORD',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
        tags: {
            name: 'login',
        }
    };

    return http.post(url, payload, params); 
}