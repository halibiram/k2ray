import re
from playwright.sync_api import sync_playwright, Page, expect

def run_verification(page: Page):
    """
    This script verifies that a user can navigate to the 2FA setup page,
    initiate the setup process, and see the QR code.
    """
    # 1. Arrange: Go to the login page.
    page.goto("http://localhost:5173/login")

    # 2. Act: Log in as the default admin user.
    page.get_by_label("Username").fill("admin")
    page.get_by_label("Password").fill("password")
    page.get_by_role("button", name="Sign in").click()

    # 3. Act: Navigate to the 2FA setup page.
    # Wait for the dashboard to confirm login was successful.
    expect(page.get_by_role("heading", name="Dashboard")).to_be_visible(timeout=10000)
    page.goto("http://localhost:5173/settings/2fa")

    # 4. Assert: Check that we are on the setup page.
    expect(page.get_by_role("heading", name="Set Up Two-Factor Authentication")).to_be_visible()

    # 5. Act: Click the button to start the 2FA setup process.
    # This might fail if the API call is not mocked, but we expect it to render the QR code UI
    # based on the frontend logic. For this test, we'll assume the API call works
    # and the frontend responds correctly.
    page.get_by_role("button", name="Enable 2FA").click()

    # 6. Assert: Wait for the QR code to appear and verify it's visible.
    qr_code_image = page.get_by_alt_text("2FA QR Code")
    expect(qr_code_image).to_be_visible()

    # 7. Screenshot: Capture the final result for visual verification.
    screenshot_path = "/app/jules-scratch/verification/2fa_setup_verification.png"
    page.screenshot(path=screenshot_path)
    print(f"Screenshot saved to {screenshot_path}")

with sync_playwright() as p:
    browser = p.chromium.launch(headless=True)
    page = browser.new_page()
    try:
        run_verification(page)
    finally:
        browser.close()