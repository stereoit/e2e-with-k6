import http from 'k6/http';
import { sleep, group, check } from 'k6';

const options = {
    vus: 1,
    iteration: 1,
    thresholds: {
        http_req_duration: ['p(99)<1500'], // 99% of requests must complete below 1.5s
    },
    // duration: '10s',
};
const SLEEP_DURATION = 0.1;

const baseUrl = __ENV.APP_HOST;

export default function () {
    var body;
    var articleID;
    var article;

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
        tags: {
            name: 'test',
        },
    };


    // const login_response = login();
    // check(login_response, {
    //     'is status 200': (r) => r.status === 200,
    //     'is access token present': (r) => r.json().hasOwnProperty('accessToken'),
    // })

    // params['Auth'] = `Bearer ${login_response.json()['accessToken']}`;
    // sleep(SLEEP_DURATION);

    group('Articles', (_) => {
        group('visit articles listing page', function () {
            params.tags.name = 'get-all-articles'
            // Get all articles
            const articles_response = http.get(
                `${baseUrl}/articles`, params
            );
            check(articles_response, {
                'is status 200': (r) => r.status === 200,
                'retrieved articles': (r) => r.json().length > 0,
            });
        });

        sleep(SLEEP_DURATION);


        // Create new article
        group('create new article', function () {
            params.tags.name = 'create-article';
            body = JSON.stringify({
                title: "test-title",
                slug: "test-slug",
            });
            const create_article_response = http.post(
                `${baseUrl}/articles`,
                body,
                params
            )
            check(create_article_response, {
                'is status 201': (r) => r.status === 201,
                'id exists': (r) => r.json("id") != null,
            });
            articleID = create_article_response.json("id");
        });
        sleep(SLEEP_DURATION);
        
        
        // Get article details
        group('Get article details', function () {
            params.tags.name = 'get-article-details';
            const resp = http.get(
                `${baseUrl}/articles/${articleID}`, params
                );
            check(resp, {
                'is status 200': (r) => r.status === 200,
                'title matches': (r) => r.json("title") == "test-title",
                'slug matches': (r) => r.json("slug") == "test-slug",
            });
            article = resp.json();
        });
        sleep(SLEEP_DURATION);

        // // Update Article
        // TODO - this one is curently broken
        // group('update article', function () {
        //     params.tags.name = 'update-article';
        //     article.slug = "test-slug-for-update";
        //     const resp = http.put(
        //         `${baseUrl}/articles/${articleID}`, 
        //         article,
        //         params
        //     );
        //     check(resp, {
        //         'is status 200': (r) => r.status === 200,
        //         'title matches': (r) => r.json("title") == "test-title",
        //         'slug matches': (r) => r.json("slug") == "test-slug",
        //     });
        // });
        // sleep(SLEEP_DURATION);

        // Delete article
        group('delete article', function () {
            params.tags.name = 'delete-article';
            const resp = http.del(
                `${baseUrl}/articles/${articleID}`,
                params
            );
            check(resp, {
                'is status 200': (r) => r.status === 200,
            });
        });
        sleep(SLEEP_DURATION);
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
            // TODO add the required header usin github secrets
        },
        tags: {
            name: 'login',
        }
    };

    return http.post(url, payload, params); 
}