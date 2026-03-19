import { test, expect } from '../support/fixtures';

test.describe('FAQ Plugin Smoke', () => {
  test('should load FAQ page', async ({ page }) => {
    await page.goto('/faq');
    await expect(page.getByRole('heading', { name: /faq/i })).toBeVisible();
  });

  test('should seed FAQ entry via API', async ({ faqFactory }) => {
    const entry = await faqFactory.seedFaq();
    expect(entry.question).toBeTruthy();
  });
});
