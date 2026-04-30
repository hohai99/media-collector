import { test, expect } from './fixtures';

test.describe('Home Page — Media Grid', () => {
  test('should display media cards from all collections', async ({ page }) => {
    await page.goto('/');
    // Wait for media cards to render
    await page.waitForSelector('.media-card', { timeout: 5000 });
    const cards = page.locator('.media-card');
    await expect(cards).toHaveCount(5); // 5 mock media items
  });

  test('should show correct media names', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.media-card');
    await expect(page.locator('.card-name').first()).toBeVisible();
  });

  test('should display type badges on media cards', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.media-card');
    const imageBadges = page.locator('.badge-image');
    const videoBadges = page.locator('.badge-video');
    const audioBadges = page.locator('.badge-audio');
    // Our mock data has 3 images, 1 video, 1 audio
    await expect(imageBadges).toHaveCount(3);
    await expect(videoBadges).toHaveCount(1);
    await expect(audioBadges).toHaveCount(1);
  });

  test('should select/deselect a media card on click', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.media-card');
    const firstCard = page.locator('.media-card').first();

    // Click to select
    await firstCard.click();
    await expect(firstCard).toHaveClass(/selected/);

    // Click again to deselect
    await firstCard.click();
    await expect(firstCard).not.toHaveClass(/selected/);
  });

  test('should show selection count when items are selected', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.media-card');

    // Select two cards
    await page.locator('.media-card').nth(0).click();
    await page.locator('.media-card').nth(1).click();

    await expect(page.locator('text=2 selected')).toBeVisible();
  });

  test('should clear selection when Clear button is clicked', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.media-card');

    await page.locator('.media-card').nth(0).click();
    await expect(page.locator('text=1 selected')).toBeVisible();

    await page.click('button:has-text("Clear")');
    await expect(page.locator('text=1 selected')).not.toBeVisible();
  });

  test('should open move modal when Move button is clicked', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.media-card');

    // Select a card
    await page.locator('.media-card').nth(0).click();
    await page.click('button:has-text("Move")');

    await expect(page.locator('.modal-overlay')).toBeVisible();
    await expect(page.locator('text=Select target collection')).toBeVisible();
  });

  test('should close move modal on Cancel', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.media-card');

    await page.locator('.media-card').nth(0).click();
    await page.click('button:has-text("Move")');
    await expect(page.locator('.modal-overlay')).toBeVisible();

    await page.click('button:has-text("Cancel")');
    await expect(page.locator('.modal-overlay')).not.toBeVisible();
  });
});
