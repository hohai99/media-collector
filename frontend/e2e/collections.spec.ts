import { test, expect } from './fixtures';

test.describe('Sidebar — Collections', () => {
  test('should display collections in sidebar', async ({ page }) => {
    await page.goto('/');
    // Collections load asynchronously after mount — use expect with auto-retry
    await expect(page.locator('.collection-item').filter({ hasText: 'Landscapes' })).toBeVisible({ timeout: 5000 });
    await expect(page.locator('.collection-item').filter({ hasText: 'Portraits' })).toBeVisible({ timeout: 5000 });
  });

  test('should show new collection input when + is clicked', async ({ page }) => {
    await page.goto('/');
    // Click the + button in the Collections section header
    await page.click('.section-header button');
    await expect(page.locator('.new-collection-form input')).toBeVisible();
  });

  test('should create a new collection', async ({ page }) => {
    await page.goto('/');
    await page.click('.section-header button');

    const input = page.locator('.new-collection-form input');
    await input.fill('Vacation Photos');
    await input.press('Enter');

    // Toast should confirm creation
    await expect(page.locator('.toast.success')).toBeVisible({ timeout: 5000 });
  });

  test('should select a collection and update header', async ({ page }) => {
    await page.goto('/');
    // Wait for collections to appear
    const landscapeItem = page.locator('.collection-item').filter({ hasText: 'Landscapes' });
    await expect(landscapeItem).toBeVisible({ timeout: 5000 });

    await landscapeItem.click();

    // Header should show collection name
    await expect(page.locator('h1')).toContainText('Landscapes', { timeout: 5000 });
  });

  test('should show Back to All Media button when collection is selected', async ({ page }) => {
    await page.goto('/');
    const landscapeItem = page.locator('.collection-item').filter({ hasText: 'Landscapes' });
    await expect(landscapeItem).toBeVisible({ timeout: 5000 });
    await landscapeItem.click();

    await expect(page.locator('button').filter({ hasText: 'All Media' })).toBeVisible({ timeout: 5000 });
  });

  test('should return to all media view', async ({ page }) => {
    await page.goto('/');
    const portraitItem = page.locator('.collection-item').filter({ hasText: 'Portraits' });
    await expect(portraitItem).toBeVisible({ timeout: 5000 });
    await portraitItem.click();
    await expect(page.locator('h1')).toContainText('Portraits', { timeout: 5000 });

    await page.locator('button').filter({ hasText: 'All Media' }).click();
    await expect(page.locator('h1')).toContainText('All Media', { timeout: 5000 });
  });
});
