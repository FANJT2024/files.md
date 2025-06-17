const { test, expect } = require('@playwright/test');

test.describe('Files.md Text Editor Sync Tests', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('http://app.localhost:8080/');

        await page.waitForSelector('.CodeMirror', { timeout: 10000 });
        await page.waitForSelector('#sidebar-tree', { timeout: 5000 });
    });

    test('should load the Files.md editor', async ({ page }) => {
        await expect(page).toHaveTitle('Files.md (Alpha version)');

        await expect(page.locator('#sidebar')).toBeVisible();
        await expect(page.locator('.CodeMirror')).toBeVisible();
        await expect(page.locator('#open-folder')).toBeVisible();
    },)
});