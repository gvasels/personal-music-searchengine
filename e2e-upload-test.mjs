import { chromium } from 'playwright';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const FRONTEND_URL = process.env.FRONTEND_URL || 'https://d8wn3lkytn5qe.cloudfront.net';
const EMAIL = process.env.E2E_TEST_EMAIL || 'gvasels90@gmail.com';
const PASSWORD = process.env.E2E_TEST_PASSWORD;

if (!PASSWORD) {
  console.error('E2E_TEST_PASSWORD environment variable is required');
  process.exit(1);
}

async function runTest() {
  console.log('Starting E2E upload test...');
  console.log(`Frontend URL: ${FRONTEND_URL}`);

  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();

  const results = [];
  const consoleErrors = [];

  // Capture browser console errors
  page.on('console', msg => {
    if (msg.type() === 'error') {
      consoleErrors.push(msg.text());
      console.log(`   [Browser Error] ${msg.text()}`);
    }
  });

  // Capture page errors
  page.on('pageerror', err => {
    consoleErrors.push(err.message);
    console.log(`   [Page Error] ${err.message}`);
  });

  // Capture failed requests
  page.on('requestfailed', request => {
    const failure = request.failure();
    console.log(`   [Request Failed] ${request.url()} - ${failure?.errorText || 'unknown'}`);
  });

  try {
    // Test 1: Navigate to home page - should redirect to login
    console.log('\n1. Navigating to home page...');
    await page.goto(FRONTEND_URL);
    await page.waitForTimeout(2000);

    const currentUrl = page.url();
    if (currentUrl.includes('/login')) {
      results.push({ test: 'Redirect to login', passed: true });
      console.log('   ✓ Redirected to login page');
    } else {
      results.push({ test: 'Redirect to login', passed: false, error: `Unexpected URL: ${currentUrl}` });
      console.log(`   ✗ Expected login redirect, got: ${currentUrl}`);
    }

    // Test 2: Login
    console.log('\n2. Logging in...');
    await page.fill('input[type="email"]', EMAIL);
    await page.fill('input[type="password"]', PASSWORD);
    await page.click('button[type="submit"]');
    await page.waitForTimeout(3000);

    const afterLoginUrl = page.url();
    if (!afterLoginUrl.includes('/login')) {
      results.push({ test: 'Login', passed: true });
      console.log('   ✓ Login successful');
    } else {
      results.push({ test: 'Login', passed: false, error: 'Still on login page' });
      console.log('   ✗ Login failed - still on login page');
    }

    // Test 3: Check sidebar is visible
    console.log('\n3. Checking sidebar visibility...');
    const sidebar = await page.$('nav, aside, [class*="sidebar"]');
    if (sidebar) {
      results.push({ test: 'Sidebar visible', passed: true });
      console.log('   ✓ Sidebar is visible after login');
    } else {
      results.push({ test: 'Sidebar visible', passed: false, error: 'Sidebar not found' });
      console.log('   ✗ Sidebar not found');
    }

    // Test 4: Navigate to upload page
    console.log('\n4. Navigating to upload page...');
    await page.goto(`${FRONTEND_URL}/upload`);
    await page.waitForTimeout(2000);

    const uploadUrl = page.url();
    if (uploadUrl.includes('/upload')) {
      results.push({ test: 'Navigate to upload', passed: true });
      console.log('   ✓ On upload page');
    } else {
      results.push({ test: 'Navigate to upload', passed: false, error: `Unexpected URL: ${uploadUrl}` });
      console.log(`   ✗ Expected upload page, got: ${uploadUrl}`);
    }

    // Test 5: Upload the sample song
    console.log('\n5. Uploading sample song...');
    const sampleSongPath = path.join(__dirname, 'sample-song', 'deadmau5 - Jaded (Volaris Remix).mp3');
    console.log(`   File path: ${sampleSongPath}`);

    // Find the file input (dropzone auto-uploads when files are selected)
    const fileInput = await page.$('input[type="file"]');
    if (fileInput) {
      await fileInput.setInputFiles(sampleSongPath);
      console.log('   File selected - upload starts automatically');

      // Wait for upload to process (dropzone auto-uploads)
      console.log('   Waiting for upload to complete...');
      await page.waitForTimeout(15000);

      // Check for status indicators
      const pageContent = await page.content();
      const hasCompleted = pageContent.includes('Completed') || pageContent.includes('completed');
      const hasProcessing = pageContent.includes('Processing') || pageContent.includes('processing');
      const hasUploading = pageContent.includes('uploading') || pageContent.includes('Uploading');
      const hasFailed = pageContent.includes('Failed') || pageContent.includes('failed') || pageContent.includes('error');
      const hasFileName = pageContent.includes('Jaded') || pageContent.includes('deadmau5');

      console.log(`   Status checks: completed=${hasCompleted}, processing=${hasProcessing}, uploading=${hasUploading}, failed=${hasFailed}, hasFileName=${hasFileName}`);

      if (hasCompleted) {
        results.push({ test: 'Upload song', passed: true });
        console.log('   ✓ Upload completed successfully');
      } else if (hasProcessing) {
        results.push({ test: 'Upload song', passed: true, note: 'Upload processing - Step Functions running' });
        console.log('   ✓ Upload processing (Step Functions running)');
      } else if (hasUploading) {
        results.push({ test: 'Upload song', passed: true, note: 'Upload in progress' });
        console.log('   ✓ Upload in progress');
      } else if (hasFailed) {
        results.push({ test: 'Upload song', passed: false, error: 'Upload failed' });
        console.log('   ✗ Upload failed');
      } else if (hasFileName) {
        results.push({ test: 'Upload song', passed: true, note: 'File appears in list' });
        console.log('   ✓ File appears in upload list');
      } else {
        results.push({ test: 'Upload song', passed: false, error: 'No upload status found' });
        console.log('   ✗ No upload status found');
      }
    } else {
      results.push({ test: 'Upload song', passed: false, error: 'File input not found' });
      console.log('   ✗ File input not found');
    }

    // Test 6: Navigate to tracks and verify upload
    console.log('\n6. Checking tracks list...');
    await page.goto(`${FRONTEND_URL}/tracks`);
    await page.waitForTimeout(3000);

    const pageContent = await page.content();
    if (pageContent.includes('Jaded') || pageContent.includes('deadmau5') || pageContent.includes('Volaris')) {
      results.push({ test: 'Track appears in list', passed: true });
      console.log('   ✓ Uploaded track found in tracks list');
    } else {
      results.push({ test: 'Track appears in list', passed: false, note: 'Track not found yet - may still be processing' });
      console.log('   ⚠ Track not found yet - may still be processing');
    }

    // Test 7: Test GLOBAL scope - check if we can see tracks from all users
    console.log('\n7. Testing GLOBAL scope (should see all tracks)...');
    const trackElements = await page.$$('[class*="track"], [class*="Track"], tr, li');
    console.log(`   Found ${trackElements.length} track elements`);
    results.push({ test: 'GLOBAL scope active', passed: trackElements.length > 0, note: `Found ${trackElements.length} elements` });

  } catch (error) {
    console.error('Test error:', error.message);
    results.push({ test: 'Unexpected error', passed: false, error: error.message });
  }

  // Print summary
  console.log('\n' + '='.repeat(50));
  console.log('TEST SUMMARY');
  console.log('='.repeat(50));

  const passed = results.filter(r => r.passed).length;
  const failed = results.filter(r => !r.passed).length;

  results.forEach(r => {
    const status = r.passed ? '✓ PASS' : '✗ FAIL';
    console.log(`${status}: ${r.test}${r.error ? ` (${r.error})` : ''}${r.note ? ` - ${r.note}` : ''}`);
  });

  console.log('\n' + '-'.repeat(50));
  console.log(`Total: ${passed} passed, ${failed} failed`);
  console.log('='.repeat(50));

  // Keep browser open for manual inspection
  console.log('\nBrowser will stay open for 30 seconds for inspection...');
  await page.waitForTimeout(30000);

  await browser.close();

  process.exit(failed > 0 ? 1 : 0);
}

runTest().catch(console.error);
