import { test, expect } from './fixtures';

test.describe('Import Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.click('.nav-item:has-text("Import")');
  });

  test('should show the Import page heading', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Import Media');
  });

  test('should show Change Folder button when master folder is set', async ({ page }) => {
    // Our mock pre-sets masterFolder, so it should say "Change Folder"
    await expect(page.locator('button').filter({ hasText: 'Change Folder' })).toBeVisible({ timeout: 5000 });
  });

  test('should show Ready to Scan state', async ({ page }) => {
    // Since master folder is pre-set, it should show the ready-to-scan state
    await expect(page.locator('text=Ready to Scan')).toBeVisible({ timeout: 5000 });
  });

  test('should scan and display media summary', async ({ page }) => {
    // Scan
    await page.click('button:has-text("Scan Media")');

    // Summary should show counts — wait for the summary grid to appear
    await expect(page.locator('.summary-number').first()).toBeVisible({ timeout: 5000 });
    // Should show Total Files count (5 items in mock)
    await expect(page.locator('.summary-number').first()).toHaveText('5');
  });

  test('should filter media by type after scan', async ({ page }) => {
    await page.click('button:has-text("Scan Media")');
    await expect(page.locator('.summary-number').first()).toBeVisible({ timeout: 5000 });

    // Click Images filter
    await page.click('button:has-text("Images")');
    await expect(page.locator('.media-card')).toHaveCount(3, { timeout: 5000 }); // 3 images in mock

    // Click Videos filter
    await page.click('button:has-text("Videos")');
    await expect(page.locator('.media-card')).toHaveCount(1, { timeout: 5000 });

    // Click Audio filter
    await page.click('button:has-text("Audio")');
    await expect(page.locator('.media-card')).toHaveCount(1, { timeout: 5000 });

    // Click All filter
    await page.click('button:has-text("All")');
    await expect(page.locator('.media-card')).toHaveCount(5, { timeout: 5000 });
  });
});
