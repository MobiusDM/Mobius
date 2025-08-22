import '@testing-library/jest-dom';
import { vi } from 'vitest';

// Mock SvelteKit modules globally
vi.mock('$app/environment', () => ({
  browser: false,
  dev: true,
  building: false
}));

vi.mock('$app/stores', () => ({
  page: {
    subscribe: (callback: any) => {
      callback({ url: { pathname: '/' } });
      return () => {};
    }
  },
  navigating: {
    subscribe: (callback: any) => {
      callback(null);
      return () => {};
    }
  }
}));

vi.mock('$app/navigation', () => ({
  goto: vi.fn(),
  invalidate: vi.fn(),
  invalidateAll: vi.fn(),
  preloadData: vi.fn(),
  preloadCode: vi.fn(),
  beforeNavigate: vi.fn(),
  afterNavigate: vi.fn()
}));
