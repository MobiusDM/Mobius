import { describe, it, expect, beforeEach, vi } from 'vitest';
import APIClient from '$lib/api';
import axios from 'axios';
import { localStorageMock } from './test-setup';

// Create complete axios mock
const mockAxiosInstance = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
  interceptors: {
    request: {
      use: vi.fn()
    },
    response: {
      use: vi.fn()
    }
  }
};

// Mock axios
vi.mock('axios', () => ({
  default: {
    create: vi.fn(() => mockAxiosInstance)
  }
}));

const mockedAxios = vi.mocked(axios, true);
mockedAxios.create = vi.fn(() => mockAxiosInstance as any);

describe('API Client', () => {
  let apiClient: APIClient;

  beforeEach(() => {
    // Clear only the mock instance methods, not the axios.create mock
    mockAxiosInstance.get.mockClear();
    mockAxiosInstance.post.mockClear();
    mockAxiosInstance.put.mockClear();
    mockAxiosInstance.delete.mockClear();
    localStorageMock.getItem.mockClear();
    localStorageMock.setItem.mockClear();
    localStorageMock.removeItem.mockClear();
    
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

      mockAxiosInstance.post.mockResolvedValueOnce({
        data: mockResponse,
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {}
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

      mockAxiosInstance.get.mockResolvedValueOnce({
        data: mockDevices,
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {}
      });

      const result = await apiClient.getDevices({
        limit: 10,
        platform: 'windows'
      });

      expect(mockAxiosInstance.get).toHaveBeenCalledWith(
        '/devices',
        expect.objectContaining({
          params: {
            limit: 10,
            platform: 'windows'
          }
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

      mockAxiosInstance.get.mockRejectedValueOnce({
        response: {
          status: 401,
          data: { error: 'Unauthorized' }
        }
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
