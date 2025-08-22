import { describe, it, expect, beforeEach, vi } from 'vitest';
import APIClient from '$lib/api';

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
};
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

// Mock fetch
global.fetch = vi.fn();

describe('API Client', () => {
  let apiClient: APIClient;

  beforeEach(() => {
    vi.clearAllMocks();
    apiClient = new APIClient('http://localhost:8081/api/v1');
  });

  describe('Authentication', () => {
    it('should store token after successful login', async () => {
      const mockResponse = {
        token: 'test-token',
        user: {
          id: '1',
          email: 'admin@mobius.local',
          name: 'Admin User',
          role: 'admin'
        }
      };

      (global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse)
      });

      const result = await apiClient.login({
        email: 'admin@mobius.local',
        password: 'admin123'
      });

      expect(localStorageMock.setItem).toHaveBeenCalledWith('auth_token', 'test-token');
      expect(result).toEqual(mockResponse);
    });

    it('should return true for authenticated user', () => {
      localStorageMock.getItem.mockReturnValue('test-token');
      
      expect(apiClient.isAuthenticated()).toBe(true);
    });

    it('should return false for unauthenticated user', () => {
      localStorageMock.getItem.mockReturnValue(null);
      
      expect(apiClient.isAuthenticated()).toBe(false);
    });

    it('should clear token on logout', async () => {
      await apiClient.logout();
      
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth_token');
    });
  });

  describe('Device Management', () => {
    beforeEach(() => {
      localStorageMock.getItem.mockReturnValue('test-token');
    });

    it('should fetch devices with proper parameters', async () => {
      const mockDevices = {
        devices: [
          {
            id: '1',
            uuid: 'device-uuid-1',
            hostname: 'test-device',
            platform: 'windows',
            status: 'online'
          }
        ],
        total: 1
      };

      (global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockDevices)
      });

      const result = await apiClient.getDevices({
        limit: 10,
        platform: 'windows'
      });

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/devices'),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Authorization': 'Bearer test-token'
          })
        })
      );
      expect(result).toEqual(mockDevices);
    });
  });

  describe('Error Handling', () => {
    it('should handle 401 errors by clearing auth token', async () => {
      localStorageMock.getItem.mockReturnValue('test-token');
      
      // Mock window.location
      delete (window as any).location;
      window.location = { href: '' } as any;

      (global.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: () => Promise.resolve({ error: 'Unauthorized' })
      });

      try {
        await apiClient.getDevices();
      } catch (error) {
        // Error is expected
      }

      expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth_token');
    });
  });
});
