import { test, expect } from './fixtures';

test.describe('App Layout & Navigation', () => {
  test('should render the sidebar with brand name', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.brand-text')).toHaveText('Media Collector');
  });

  test('should show navigation items in sidebar', async ({ page }) => {
    await page.goto('/');
    const navItems = page.locator('.nav-item');
    await expect(navItems).toHaveCount(3);
    await expect(navItems.nth(0)).toContainText('Collections');
    await expect(navItems.nth(1)).toContainText('Import');
    await expect(navItems.nth(2)).toContainText('Player');
  });

  test('should navigate to Import page', async ({ page }) => {
    await page.goto('/');
    await page.click('.nav-item:has-text("Import")');
    await expect(page.locator('h1')).toContainText('Import Media');
  });

  test('should navigate to Player page', async ({ page }) => {
    await page.goto('/');
    await page.click('.nav-item:has-text("Player")');
    await expect(page.locator('h1')).toContainText('Player Configurations');
  });

  test('should navigate back to Collections page', async ({ page }) => {
    await page.goto('/');
    await page.click('.nav-item:has-text("Import")');
    await page.click('.nav-item:has-text("Collections")');
    await expect(page.locator('h1')).toContainText('All Media');
  });

  test('should highlight the active nav item', async ({ page }) => {
    await page.goto('/');
    const collectionsNav = page.locator('.nav-item:has-text("Collections")');
    await expect(collectionsNav).toHaveClass(/active/);

    await page.click('.nav-item:has-text("Import")');
    const importNav = page.locator('.nav-item:has-text("Import")');
    await expect(importNav).toHaveClass(/active/);
    await expect(collectionsNav).not.toHaveClass(/active/);
  });
});
