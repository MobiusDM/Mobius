import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/svelte';
import Layout from '$lib/Layout.svelte';

// Mock the API client
vi.mock('$lib/api', () => ({
  apiClient: {
    isAuthenticated: () => true
  }
}));

// Mock SvelteKit modules
vi.mock('$app/stores', () => ({
  page: {
    subscribe: (callback: any) => {
      callback({ url: { pathname: '/' } });
      return () => {};
    }
  }
}));

vi.mock('$app/navigation', () => ({
  goto: vi.fn()
}));

describe('Layout Component', () => {
  it('renders the main navigation', () => {
    render(Layout, {
      props: {}
    });

    expect(screen.getByText('Mobius MDM')).toBeInTheDocument();
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
    expect(screen.getByText('Devices')).toBeInTheDocument();
    expect(screen.getByText('Policies')).toBeInTheDocument();
    expect(screen.getByText('Applications')).toBeInTheDocument();
    expect(screen.getByText('Groups')).toBeInTheDocument();
  });

  it('has proper navigation structure', () => {
    render(Layout, {
      props: {}
    });

    const dashboardLink = screen.getByRole('link', { name: /dashboard/i });
    expect(dashboardLink).toHaveAttribute('href', '/');

    const devicesLink = screen.getByRole('link', { name: /devices/i });
    expect(devicesLink).toHaveAttribute('href', '/devices');
  });
});
