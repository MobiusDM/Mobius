# SECURITY.md Automation

This document explains the automated security documentation maintenance for Mobius.

## Overview

The `SECURITY.md` file in the repository root provides comprehensive security information including:

- Vulnerability reporting procedures
- Security features and best practices  
- Supported versions and contact information
- Deployment and operational security guidance

## Automation

### Workflow: `.github/workflows/update-security.yml`

The automation workflow runs:

- **Monthly**: 1st of each month at 2 AM UTC
- **On push**: When SECURITY.md or workflow files change
- **Manual**: Can be triggered via GitHub Actions UI

### What It Does

1. **Updates timestamps**: Keeps "Last Updated" current
2. **Increments version**: Updates security policy version  
3. **Checks vulnerabilities**: Scans npm and Go dependencies
4. **Validates format**: Ensures all required sections exist
5. **Creates issues**: Auto-creates issues for security attention

### Outputs

- **Commits**: Updates SECURITY.md with current information
- **Issues**: Creates security alert issues when vulnerabilities found
- **Validation**: Ensures documentation quality and completeness

## Maintenance

### Required Manual Updates

1. **Security Contact Email**: Update `security@domain.local` placeholder
2. **Organization Details**: Customize contact information
3. **Environment Specifics**: Adjust configuration examples for your setup

### Customization

The automation can be customized by modifying:

- **Schedule**: Change the cron schedule in the workflow
- **Content checks**: Update validation rules in the Python script
- **Issue creation**: Modify alert thresholds and templates

### Testing

Run the test script to validate automation:

```bash
# Test all automation components
/tmp/security-test/test_security_automation.sh
```

### Troubleshooting

**Common Issues:**

- **Placeholder emails**: Update security contact information
- **Missing sections**: Ensure all required sections exist in SECURITY.md
- **Workflow failures**: Check GitHub Actions logs for specific errors

**Manual workflow trigger:**

1. Go to GitHub Actions tab
2. Select "Update Security Documentation"  
3. Click "Run workflow"

## Integration

The SECURITY.md is integrated with:

- **README.md**: Links to security policy
- **GitHub Security**: Referenced in vulnerability reporting
- **Dependencies**: Connected to Dependabot and security scanning

## Best Practices

1. **Review monthly**: Check automated updates for accuracy
2. **Update contacts**: Keep security contact information current
3. **Test procedures**: Verify vulnerability reporting process works
4. **Monitor issues**: Address security alerts promptly

---

*This automation ensures Mobius maintains up-to-date security documentation following industry best practices.*