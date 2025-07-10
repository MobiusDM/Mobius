package installer_cache

import (
	"log"
	"os"
	"sync"

	"github.com/notawar/mobius/backend/cmd/osquery-perf/osquery_perf"
	"github.com/notawar/mobius/backend/pkg/file"
	"github.com/notawar/mobius/backend/server/mobius"
	"github.com/notawar/mobius/backend/server/service"
)

// Metadata holds the metadata for software installers.
// To extract the metadata, we must download the file. Once the file has been downloaded once and analyzed,
// the other agents can use the cache to get the appropriate metadata.
type Metadata struct {
	mu    sync.Mutex
	cache map[uint]*file.InstallerMetadata
	Stats *osquery_perf.Stats
}

func (c *Metadata) Get(installer *mobius.SoftwareInstallDetails, orbitClient *service.OrbitClient) (meta *file.InstallerMetadata,
	cacheMiss bool, err error,
) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		c.cache = make(map[uint]*file.InstallerMetadata, 1)
	}

	meta, ok := c.cache[installer.InstallerID]
	if !ok {
		var err error
		meta, err = c.populateMetadata(installer, orbitClient)
		if err != nil {
			return nil, false, err
		}
		c.cache[installer.InstallerID] = meta
		cacheMiss = true
	}
	return meta, cacheMiss, nil
}

func (c *Metadata) populateMetadata(installer *mobius.SoftwareInstallDetails, orbitClient *service.OrbitClient) (*file.InstallerMetadata,
	error,
) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		c.Stats.IncrementOrbitErrors()
		log.Println("level=error, create temp dir:", err)
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	var path string
	if installer.SoftwareInstallerURL != nil {
		path, err = orbitClient.DownloadSoftwareInstallerFromURL(installer.SoftwareInstallerURL.URL,
			installer.SoftwareInstallerURL.Filename, tmpDir, func(n int) {
			})
		if err != nil {
			log.Printf("level=error, msg=download software installer from URL; is CloudFront CDN configured correctly?, err=%s", err)
			c.Stats.IncrementOrbitErrors()
			return nil, err
		}
	}

	if path == "" {
		path, err = orbitClient.DownloadSoftwareInstaller(installer.InstallerID, tmpDir, func(n int) {
		})
		if err != nil {
			log.Printf("level=error, msg=download software installer, err=%s", err)
			c.Stats.IncrementOrbitErrors()
			return nil, err
		}
	}
	// Figure out what we're actually installing here and add it to software inventory
	tfr, err := mobius.NewKeepFileReader(path)
	if err != nil {
		c.Stats.IncrementOrbitErrors()
		log.Println("level=error, open installer:", err)
		return nil, err
	}
	defer tfr.Close()
	item, err := file.ExtractInstallerMetadata(tfr)
	if err != nil {
		c.Stats.IncrementOrbitErrors()
		log.Println("level=error, extract installer metadata:", err)
		return nil, err
	}
	return item, nil
}
