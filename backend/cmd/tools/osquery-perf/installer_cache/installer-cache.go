package installer_cache

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/notawar/mobius/cmd/tools/osquery-perf/osquery_perf"
	"github.com/notawar/mobius/internal/server/mobius"
	"github.com/notawar/mobius/internal/server/service"
	"github.com/notawar/mobius/pkg/file"
)

// Metadata holds the metadata for software installers.
// To extract the metadata, we must download the file. Once the file has been downloaded once and analyzed,
// the other agents can use the cache to get the appropriate metadata.
type Metadata struct {
	mu    sync.Mutex
	cache map[uint]*file.InstallerMetadata
	Stats *osquery_perf.Stats
}

func (c *Metadata) Get(installer *mobius.SoftwareInstallDetails, client *service.Client) (meta *file.InstallerMetadata,
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
		meta, err = c.populateMetadata(installer, client)
		if err != nil {
			return nil, false, err
		}
		c.cache[installer.InstallerID] = meta
		cacheMiss = true
	}
	return meta, cacheMiss, nil
}

func (c *Metadata) populateMetadata(installer *mobius.SoftwareInstallDetails, _ *service.Client) (*file.InstallerMetadata,
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
		// Download file from URL directly
		path = filepath.Join(tmpDir, installer.SoftwareInstallerURL.Filename)
		err = downloadFileFromURL(installer.SoftwareInstallerURL.URL, path)
		if err != nil {
			log.Printf("level=error, msg=download software installer from URL; err=%s", err)
			c.Stats.IncrementOrbitErrors()
			return nil, err
		}
	} else {
		// Use client API to download software installer
		// Note: This assumes client has an appropriate method or endpoint for this
		// Use a simple API call to get the file data
		path = filepath.Join(tmpDir, fmt.Sprintf("installer-%d", installer.InstallerID))
		log.Printf("level=info, msg=downloading software installer %d to %s", installer.InstallerID, path)

		// Since direct client methods aren't available, we would need to implement
		// API calls to download the file. For now, this is a placeholder implementation
		// that needs to be completed when the API is properly defined.
		err = fmt.Errorf("downloading software installers not implemented in current client")
		log.Printf("level=error, msg=download software installer, err=%s", err)
		c.Stats.IncrementOrbitErrors()
		return nil, err
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

// Helper function to download a file from URL to a specific path
func downloadFileFromURL(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
