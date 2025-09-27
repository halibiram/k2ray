import { test, expect } from '@playwright/test';

test.describe('App E2E', () => {
  test('should navigate to the login page and match the screenshot', async ({ page }) => {
    // Navigate to the root of the app, which should redirect to /login
    await page.goto('/');

    // Wait for the login page to load by checking for a specific element
    await expect(page.locator('h2:has-text("Login to your account")')).toBeVisible();

    // Assert the page title
    await expect(page).toHaveTitle(/K2Ray/);

    // Take a screenshot for visual regression testing
    // The first run will create a 'golden' or 'snapshot' file.
    // Subsequent runs will compare against this snapshot.
    await expect(page).toHaveScreenshot('login-page.png', { fullPage: true });
  });
});