import { test, expect } from './fixtures';

test.describe('Player Builder Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.click('.nav-item:has-text("Player")');
  });

  test('should show the Player Configurations heading', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Player Configurations');
  });

  test('should show the create form', async ({ page }) => {
    await expect(page.locator('text=Create New Config')).toBeVisible();
    await expect(page.locator('#config-name')).toBeVisible();
    await expect(page.locator('#total-time')).toBeVisible();
  });

  test('should show existing player configs', async ({ page }) => {
    // Mock has 1 config: "Morning Slideshow"
    await expect(page.locator('.config-card').filter({ hasText: 'Morning Slideshow' })).toBeVisible({ timeout: 5000 });
  });

  test('should show Play and Delete buttons on saved configs', async ({ page }) => {
    await expect(page.locator('.config-card').first()).toBeVisible({ timeout: 5000 });
    // Play button inside config-card
    await expect(page.locator('.config-card button').filter({ hasText: 'Play' })).toBeVisible();
    // Delete button (🗑) inside config-card
    await expect(page.locator('.config-card .btn-danger')).toBeVisible();
  });

  test('should fill in config name', async ({ page }) => {
    await page.fill('#config-name', 'Evening Slideshow');
    await expect(page.locator('#config-name')).toHaveValue('Evening Slideshow');
  });

  test('should set total time', async ({ page }) => {
    await page.fill('#total-time', '10');
    await expect(page.locator('#total-time')).toHaveValue('10');
  });

  test('should show collection checkboxes', async ({ page }) => {
    // Wait for collections to load in the folder tree
    const checkboxes = page.locator('.folder-tree input[type="checkbox"]');
    await expect(checkboxes.first()).toBeVisible({ timeout: 5000 });
  });

  test('should delete a player config', async ({ page }) => {
    await expect(page.locator('.config-card').first()).toBeVisible({ timeout: 5000 });

    // Click delete
    await page.locator('.config-card .btn-danger').first().click();

    // Toast should confirm deletion
    await expect(page.locator('.toast.success')).toBeVisible({ timeout: 5000 });
  });
});

test.describe('Player View', () => {
  test('should navigate to player when Play is clicked', async ({ page }) => {
    await page.goto('/');
    await page.click('.nav-item:has-text("Player")');

    // Wait for config card
    await expect(page.locator('.config-card').first()).toBeVisible({ timeout: 5000 });

    // Click Play on the existing config
    await page.locator('.config-card button').filter({ hasText: 'Play' }).first().click();

    // Player should show — either the fullscreen view or empty state
    // The mock resolves media for col-1 which has sunset.jpg
    const playerOrEmpty = page.locator('.player-fullscreen, :text("No Media to Play")');
    await expect(playerOrEmpty.first()).toBeVisible({ timeout: 5000 });
  });
});
