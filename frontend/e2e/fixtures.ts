/**
 * E2E test fixtures for Media Collector.
 *
 * Since Playwright tests run in a real browser without the Wails runtime,
 * Go bindings (window.go.main.App.*) are not available. These fixtures
 * inject mock implementations before each test so the Svelte app can
 * render and respond to user interactions.
 */
import { test as base, expect, type Page } from '@playwright/test';

// Sample mock data
export const mockCollections = [
  { id: 'col-1', name: 'Landscapes', path: 'C:/media/Landscapes', parentId: '' },
  { id: 'col-2', name: 'Portraits', path: 'C:/media/Portraits', parentId: '' },
  { id: 'col-3', name: 'Mountains', path: 'C:/media/Landscapes/Mountains', parentId: 'col-1' },
];

export const mockMedia = [
  { id: 'm-1', name: 'sunset.jpg', path: 'C:/media/sunset.jpg', type: 'image', size: 2048000, createdAt: '2025-01-15T10:30:00Z', collectionId: 'col-1' },
  { id: 'm-2', name: 'portrait.png', path: 'C:/media/portrait.png', type: 'image', size: 1536000, createdAt: '2025-02-20T14:00:00Z', collectionId: 'col-2' },
  { id: 'm-3', name: 'clip.mp4', path: 'C:/media/clip.mp4', type: 'video', size: 51200000, createdAt: '2025-03-10T08:00:00Z', collectionId: '' },
  { id: 'm-4', name: 'song.mp3', path: 'C:/media/song.mp3', type: 'audio', size: 4096000, createdAt: '2025-04-01T12:00:00Z', collectionId: '' },
  { id: 'm-5', name: 'mountain.jpg', path: 'C:/media/Landscapes/Mountains/mountain.jpg', type: 'image', size: 3072000, createdAt: '2025-05-05T16:00:00Z', collectionId: 'col-3' },
];

export const mockPlayerConfigs = [
  { id: 'pc-1', name: 'Morning Slideshow', totalPlayTime: 180, transitionTime: 36, collections: ['col-1'], mediaIds: [] },
];

/**
 * Inject Wails Go binding mocks into the page.
 * This creates `window.go.main.App` with methods that return mock data.
 */
async function injectWailsMocks(page: Page) {
  await page.addInitScript(`
    // Stub the Wails runtime
    window.runtime = window.runtime || {};
    window.runtime.EventsOn = function() {};
    window.runtime.EventsOff = function() {};
    window.runtime.EventsOnce = function() {};
    window.runtime.EventsEmit = function() {};
    window.runtime.LogDebug = function() {};
    window.runtime.LogInfo = function() {};
    window.runtime.LogWarning = function() {};
    window.runtime.LogError = function() {};

    // Mock data store — masterFolder is pre-set so collections load on init
    window.__mockData = {
      masterFolder: 'C:/media',
      collections: ${JSON.stringify(mockCollections)},
      media: ${JSON.stringify(mockMedia)},
      playerConfigs: ${JSON.stringify(mockPlayerConfigs)},
    };

    // Stub go bindings — these are what the Wails-generated JS calls
    window.go = window.go || {};
    window.go.main = window.go.main || {};
    window.go.main.App = {
      SelectMasterFolder: async () => {
        window.__mockData.masterFolder = 'C:/media';
        return 'C:/media';
      },
      GetMasterFolder: async () => window.__mockData.masterFolder,
      ScanMedia: async () => window.__mockData.media,
      GetAllMedia: async () => window.__mockData.media,
      GetMediaByCollection: async (colId) =>
        window.__mockData.media.filter(m => m.collectionId === colId),
      GetUncategorizedMedia: async () =>
        window.__mockData.media.filter(m => !m.collectionId),
      MoveMedia: async (ids, colId) => {
        ids.forEach(id => {
          const m = window.__mockData.media.find(x => x.id === id);
          if (m) m.collectionId = colId;
        });
      },
      CreateCollection: async (name, parentId) => {
        const newCol = {
          id: 'col-' + Date.now(),
          name,
          path: 'C:/media/' + name,
          parentId: parentId || '',
        };
        window.__mockData.collections.push(newCol);
        return newCol;
      },
      GetCollectionTree: async () => window.__mockData.collections,
      DeleteCollection: async (id) => {
        window.__mockData.collections = window.__mockData.collections.filter(c => c.id !== id);
      },
      SyncCollections: async () => {},
      CreatePlayerConfig: async (input) => {
        const cfg = {
          id: 'pc-' + Date.now(),
          name: input.name,
          totalPlayTime: input.totalPlayTime,
          transitionTime: Math.floor(input.totalPlayTime / Math.max(input.collections?.length || 1, 1)),
          collections: input.collections || [],
          mediaIds: input.mediaIds || [],
        };
        window.__mockData.playerConfigs.push(cfg);
        return cfg;
      },
      GetPlayerConfigs: async () => window.__mockData.playerConfigs,
      GetPlayerConfig: async (id) => window.__mockData.playerConfigs.find(c => c.id === id) || null,
      DeletePlayerConfig: async (id) => {
        window.__mockData.playerConfigs = window.__mockData.playerConfigs.filter(c => c.id !== id);
      },
      ResolvePlayerMedia: async (configId) => {
        const cfg = window.__mockData.playerConfigs.find(c => c.id === configId);
        if (!cfg) return [];
        let media = [];
        for (const colId of (cfg.collections || [])) {
          media = media.concat(window.__mockData.media.filter(m => m.collectionId === colId));
        }
        return media;
      },
      GetAbsolutePath: async (p) => p,
    };
  `);
}

// Extended test fixture with mocks pre-injected
export const test = base.extend({
  page: async ({ page }, use) => {
    await injectWailsMocks(page);
    await use(page);
  },
});

export { expect };
