import axios from 'axios';
import assert from 'assert';

const url = 'http://localhost:5173';

async function runSmokeTest() {
  try {
    console.log(`Sending request to ${url}...`);
    const response = await axios.get(url);

    assert.strictEqual(response.status, 200, `Expected status code 200, but got ${response.status}`);
    console.log('✔ Status code is 200');

    const html = response.data;
    assert.ok(html.includes('<title>K2Ray</title>'), 'HTML response should contain the correct title tag');
    console.log('✔ HTML title is correct');

    console.log('\nSmoke test passed successfully!');
    process.exit(0);
  } catch (error) {
    console.error('\nSmoke test failed:');
    console.error(error.message);
    if (error.response) {
      console.error('Response Status:', error.response.status);
      console.error('Response Data:', error.response.data);
    }
    process.exit(1);
  }
}

runSmokeTest();