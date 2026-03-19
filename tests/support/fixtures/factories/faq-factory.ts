import { faker } from '@faker-js/faker';

type FaqEntry = {
  id?: string;
  question: string;
  answer: string;
  locale: string;
};

export class FaqFactory {
  private createdFaqIds: string[] = [];

  createFaq(overrides: Partial<FaqEntry> = {}): FaqEntry {
    return {
      question: faker.lorem.sentence(),
      answer: faker.lorem.paragraph(),
      locale: 'en',
      ...overrides,
    };
  }

  async seedFaq(overrides: Partial<FaqEntry> = {}): Promise<FaqEntry> {
    const faq = this.createFaq(overrides);

    const response = await fetch(`${process.env.API_URL}/faq`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(faq),
    });

    if (!response.ok) {
      throw new Error(`Failed to seed FAQ: ${response.status} ${await response.text()}`);
    }

    const created = (await response.json()) as FaqEntry & { id: string };
    this.createdFaqIds.push(created.id);
    return created;
  }

  async cleanup(): Promise<void> {
    for (const faqId of this.createdFaqIds) {
      await fetch(`${process.env.API_URL}/faq/${faqId}`, {
        method: 'DELETE',
      });
    }
    this.createdFaqIds = [];
  }
}
