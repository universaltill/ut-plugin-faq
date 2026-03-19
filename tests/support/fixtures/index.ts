import { test as base, expect } from '@playwright/test';
import { FaqFactory } from './factories/faq-factory';

type TestFixtures = {
  faqFactory: FaqFactory;
};

export const test = base.extend<TestFixtures>({
  faqFactory: async ({}, use) => {
    const factory = new FaqFactory();
    await use(factory);
    await factory.cleanup();
  },
});

export { expect };
