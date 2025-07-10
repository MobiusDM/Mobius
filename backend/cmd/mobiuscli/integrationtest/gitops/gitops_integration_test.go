package gitops

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/notawar/mobius/backend/cmd/mobiuscli/mobiuscli"
	"github.com/notawar/mobius/backend/cmd/mobiuscli/integrationtest"
	"github.com/notawar/mobius/backend/server/config"
	"github.com/notawar/mobius/backend/server/datastore/redis/redistest"
	"github.com/notawar/mobius/backend/server/mobius"
	appleMdm "github.com/notawar/mobius/backend/server/mdm/apple"
	"github.com/notawar/mobius/backend/server/mdm/nanodep/tokenpki"
	"github.com/notawar/mobius/backend/server/service"
	"github.com/notawar/mobius/backend/server/test"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestIntegrationsGitops(t *testing.T) {
	testingSuite := new(integrationGitopsTestSuite)
	testingSuite.WithServer.Suite = &testingSuite.Suite
	suite.Run(t, testingSuite)
}

type integrationGitopsTestSuite struct {
	suite.Suite
	integrationtest.WithServer
	mobiusCfg config.MobiusConfig
}

func (s *integrationGitopsTestSuite) SetupSuite() {
	s.WithDS.SetupSuite("integrationGitopsTestSuite")

	appConf, err := s.DS.AppConfig(context.Background())
	require.NoError(s.T(), err)
	appConf.MDM.EnabledAndConfigured = true
	appConf.MDM.AppleBMEnabledAndConfigured = true
	appConf.MDM.WindowsEnabledAndConfigured = true
	err = s.DS.SaveAppConfig(context.Background(), appConf)
	require.NoError(s.T(), err)

	testCert, testKey, err := appleMdm.NewSCEPCACertKey()
	require.NoError(s.T(), err)
	testCertPEM := tokenpki.PEMCertificate(testCert.Raw)
	testKeyPEM := tokenpki.PEMRSAPrivateKey(testKey)

	mobiusCfg := config.TestConfig()
	config.SetTestMDMConfig(s.T(), &mobiusCfg, testCertPEM, testKeyPEM, "../../../../server/service/testdata")
	mobiusCfg.Osquery.EnrollCooldown = 0

	mdmStorage, err := s.DS.NewMDMAppleMDMStorage()
	require.NoError(s.T(), err)
	depStorage, err := s.DS.NewMDMAppleDEPStorage()
	require.NoError(s.T(), err)
	scepStorage, err := s.DS.NewSCEPDepot()
	require.NoError(s.T(), err)
	redisPool := redistest.SetupRedis(s.T(), "zz", false, false, false)

	serverConfig := service.TestServerOpts{
		License: &mobius.LicenseInfo{
			Tier: mobius.TierFree,
		},
		MobiusConfig: &mobiusCfg,
		MDMStorage:  mdmStorage,
		DEPStorage:  depStorage,
		SCEPStorage: scepStorage,
		Pool:        redisPool,
		APNSTopic:   "com.apple.mgmt.External.10ac3ce5-4668-4e58-b69a-b2b5ce667589",
	}
	err = s.DS.InsertMDMConfigAssets(context.Background(), []mobius.MDMConfigAsset{
		{Name: mobius.MDMAssetSCEPChallenge, Value: []byte("scepchallenge")},
	}, nil)
	require.NoError(s.T(), err)
	users, server := service.RunServerForTestsWithDS(s.T(), s.DS, &serverConfig)
	s.T().Setenv("MOBIUS_SERVER_ADDRESS", server.URL) // mobiuscli always uses this env var in tests
	s.Server = server
	s.Users = users
	s.mobiusCfg = mobiusCfg

	appConf, err = s.DS.AppConfig(context.Background())
	require.NoError(s.T(), err)
	appConf.ServerSettings.ServerURL = server.URL
	err = s.DS.SaveAppConfig(context.Background(), appConf)
	require.NoError(s.T(), err)
}

func (s *integrationGitopsTestSuite) TearDownSuite() {
	appConf, err := s.DS.AppConfig(context.Background())
	require.NoError(s.T(), err)
	appConf.MDM.EnabledAndConfigured = false
	err = s.DS.SaveAppConfig(context.Background(), appConf)
	require.NoError(s.T(), err)
}

// TestMobiusGitops runs `mobiuscli gitops` command on configs in https://github.com/notawar/mobius-gitops repo.
// Changes to that repo may cause this test to fail.
func (s *integrationGitopsTestSuite) TestMobiusGitops() {
	t := s.T()
	const mobiusGitopsRepo = "https://github.com/notawar/mobius-gitops"

	mobiuscliConfig := s.createMobiusctlConfig()

	// Clone git repo
	repoDir := t.TempDir()
	_, err := git.PlainClone(
		repoDir, false, &git.CloneOptions{
			ReferenceName: "main",
			SingleBranch:  true,
			Depth:         1,
			URL:           mobiusGitopsRepo,
			Progress:      os.Stdout,
		},
	)
	require.NoError(t, err)

	// Set the required environment variables
	t.Setenv("MOBIUS_URL", s.Server.URL)
	t.Setenv("MOBIUS_GLOBAL_ENROLL_SECRET", "global_enroll_secret")
	globalFile := path.Join(repoDir, "default.yml")

	// Dry run
	_ = mobiuscli.RunAppForTest(t, []string{"gitops", "--config", mobiuscliConfig.Name(), "-f", globalFile, "--dry-run"})

	// Real run
	_ = mobiuscli.RunAppForTest(t, []string{"gitops", "--config", mobiuscliConfig.Name(), "-f", globalFile})

}

func (s *integrationGitopsTestSuite) createMobiusctlConfig() *os.File {
	t := s.T()
	// Create a temporary mobiuscli config file
	mobiuscliConfig, err := os.CreateTemp(t.TempDir(), "*.yml")
	require.NoError(t, err)
	// GitOps user is a premium feature, so we simply use an admin user.
	token := s.GetTestToken("admin1@example.com", test.GoodPassword)
	configStr := fmt.Sprintf(
		`
contexts:
  default:
    address: %s
    tls-skip-verify: true
    token: %s
`, s.Server.URL, token,
	)
	_, err = mobiuscliConfig.WriteString(configStr)
	require.NoError(t, err)
	return mobiuscliConfig
}

func (s *integrationGitopsTestSuite) TestMobiusGitopsWithMobiusSecrets() {
	t := s.T()
	const (
		secretName1 = "NAME"
		secretName2 = "length"
	)
	ctx := context.Background()
	mobiuscliConfig := s.createMobiusctlConfig()

	// Set the required environment variables
	t.Setenv("MOBIUS_URL", s.Server.URL)
	t.Setenv("MOBIUS_GLOBAL_ENROLL_SECRET", "global_enroll_secret")
	t.Setenv("MOBIUS_SECRET_"+secretName1, "secret_value")
	t.Setenv("MOBIUS_SECRET_"+secretName2, "2")
	globalFile := path.Join("..", "..", "mobiuscli", "testdata", "gitops", "global_integration.yml")

	// Dry run
	_ = mobiuscli.RunAppForTest(t, []string{"gitops", "--config", mobiuscliConfig.Name(), "-f", globalFile, "--dry-run"})
	secrets, err := s.DS.GetSecretVariables(ctx, []string{secretName1})
	require.NoError(t, err)
	require.Empty(t, secrets)

	// Real run
	_ = mobiuscli.RunAppForTest(t, []string{"gitops", "--config", mobiuscliConfig.Name(), "-f", globalFile})
	// Check secrets
	secrets, err = s.DS.GetSecretVariables(ctx, []string{secretName1, secretName2})
	require.NoError(t, err)
	require.Len(t, secrets, 2)
	for _, secret := range secrets {
		switch secret.Name {
		case secretName1:
			assert.Equal(t, "secret_value", secret.Value)
		case secretName2:
			assert.Equal(t, "2", secret.Value)
		default:
			t.Fatalf("unexpected secret %s", secret.Name)
		}
	}

	// Check script(s)
	scriptID, err := s.DS.GetScriptIDByName(ctx, "mobius-secret.sh", nil)
	require.NoError(t, err)
	expected, err := os.ReadFile("../../mobiuscli/testdata/gitops/lib/mobius-secret.sh")
	require.NoError(t, err)
	script, err := s.DS.GetScriptContents(ctx, scriptID)
	require.NoError(t, err)
	assert.Equal(t, expected, script)

	// Check Apple profiles
	profiles, err := s.DS.ListMDMAppleConfigProfiles(ctx, nil)
	require.NoError(t, err)
	require.Len(t, profiles, 1)
	assert.Contains(t, string(profiles[0].Mobileconfig), "$MOBIUS_SECRET_"+secretName1)
	// Check Windows profiles
	allProfiles, _, err := s.DS.ListMDMConfigProfiles(ctx, nil, mobius.ListOptions{})
	require.NoError(t, err)
	require.Len(t, allProfiles, 2)
	var windowsProfileUUID string
	for _, profile := range allProfiles {
		if profile.Platform == "windows" {
			windowsProfileUUID = profile.ProfileUUID
		}
	}
	require.NotEmpty(t, windowsProfileUUID)
	winProfile, err := s.DS.GetMDMWindowsConfigProfile(ctx, windowsProfileUUID)
	require.NoError(t, err)
	assert.Contains(t, string(winProfile.SyncML), "${MOBIUS_SECRET_"+secretName2+"}")
}
