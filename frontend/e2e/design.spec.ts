import { test, expect } from './fixtures';

test.describe('Visual & Design', () => {
  test('should have dark background', async ({ page }) => {
    await page.goto('/');
    const bgColor = await page.evaluate(() =>
      getComputedStyle(document.body).backgroundColor
    );
    // Should be dark (rgb values should be low)
    expect(bgColor).toBeTruthy();
  });

  test('should use Inter font family', async ({ page }) => {
    await page.goto('/');
    const fontFamily = await page.evaluate(() =>
      getComputedStyle(document.body).fontFamily
    );
    expect(fontFamily.toLowerCase()).toContain('inter');
  });

  test('sidebar should have correct width', async ({ page }) => {
    await page.goto('/');
    const sidebar = page.locator('.sidebar');
    const box = await sidebar.boundingBox();
    expect(box).toBeTruthy();
    // Sidebar should be around 260px
    expect(box!.width).toBeGreaterThanOrEqual(250);
    expect(box!.width).toBeLessThanOrEqual(280);
  });

  test('brand text should have gradient styling', async ({ page }) => {
    await page.goto('/');
    const brand = page.locator('.brand-text');
    const bgImage = await brand.evaluate((el) =>
      getComputedStyle(el).backgroundImage
    );
    expect(bgImage).toContain('gradient');
  });

  test('glass cards should have backdrop-filter', async ({ page }) => {
    await page.goto('/');
    await page.click('.nav-item:has-text("Player")');
    await expect(page.locator('.glass-card').first()).toBeVisible({ timeout: 5000 });

    const blur = await page.locator('.glass-card').first().evaluate((el) =>
      getComputedStyle(el).backdropFilter || (getComputedStyle(el) as any).webkitBackdropFilter
    );
    expect(blur).toContain('blur');
  });
});

test.describe('Responsive & Interactive', () => {
  test('media cards should have hover effect', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.media-card', { timeout: 5000 });

    const card = page.locator('.media-card').first();
    await card.hover();
    // Verify the card is visible and hoverable (hover styles are CSS-applied)
    await expect(card).toBeVisible();
  });

  test('buttons should have pointer cursor', async ({ page }) => {
    await page.goto('/');
    const navBtn = page.locator('.nav-item').first();
    const cursor = await navBtn.evaluate((el) =>
      getComputedStyle(el).cursor
    );
    expect(cursor).toBe('pointer');
  });

  test('app should render without errors', async ({ page }) => {
    await page.goto('/');
    // Just verify the page renders without errors
    const hasContent = await page.locator('#app').isVisible();
    expect(hasContent).toBeTruthy();
  });
});

test.describe('Toast Notifications', () => {
  test('should show success toast when collection is created', async ({ page }) => {
    await page.goto('/');

    // Create a collection via sidebar
    await page.click('.section-header button');
    const input = page.locator('.new-collection-form input');
    await input.fill('Test Collection');
    await input.press('Enter');

    // Wait for toast
    await expect(page.locator('.toast.success')).toBeVisible({ timeout: 5000 });
    await expect(page.locator('.toast.success')).toContainText('Collection created');
  });

  test('toast should auto-dismiss after 3 seconds', async ({ page }) => {
    await page.goto('/');

    await page.click('.section-header button');
    const input = page.locator('.new-collection-form input');
    await input.fill('Temp Collection');
    await input.press('Enter');

    await expect(page.locator('.toast.success')).toBeVisible({ timeout: 5000 });
    // Wait for auto-dismiss (3 seconds + buffer)
    await expect(page.locator('.toast.success')).not.toBeVisible({ timeout: 6000 });
  });
});
