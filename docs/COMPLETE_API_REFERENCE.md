# 📡 Mobius MDM Platform - Complete API Reference

## Overview
**Yes, everything in the Mobius MDM platform is API-first!** All functionality is accessible through RESTful HTTP endpoints, with optional real-time capabilities via WebSocket for live updates.

## 🔐 Authentication
All API endpoints (except health) require JWT authentication:
```bash
# Login to get JWT token
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin"}'

# Use token in subsequent requests
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8081/api/v1/devices
```

## 📋 Complete API Endpoint List

### 🏥 System Health & Status
```bash
GET  /health                           # Server health check (no auth required)
GET  /api/v1/license/status           # License information
PUT  /api/v1/license                  # Update license key
```

### 🔑 Authentication
```bash
POST /api/v1/auth/login               # User login (get JWT token)
POST /api/v1/auth/logout              # User logout
```

### 📱 Device Management
```bash
# Device Lifecycle
POST /api/v1/devices                  # Enroll new device
GET  /api/v1/devices                  # List all devices (with filtering)
GET  /api/v1/devices/{id}             # Get specific device details
PUT  /api/v1/devices/{id}             # Update device information
DELETE /api/v1/devices/{id}           # Unenroll device

# Device Operations
POST /api/v1/devices/{id}/commands    # Execute command on device
POST /api/v1/devices/{id}/osquery     # Run OSQuery on device

# Device Search & Filtering
GET  /api/v1/devices?platform=linux   # Filter by platform
GET  /api/v1/devices?status=online    # Filter by status
GET  /api/v1/devices?search=hostname  # Search by hostname/name
```

### 📋 Policy Management
```bash
# Policy Lifecycle
POST /api/v1/policies                 # Create new policy
GET  /api/v1/policies                 # List all policies
GET  /api/v1/policies/{id}            # Get specific policy
PUT  /api/v1/policies/{id}            # Update policy
DELETE /api/v1/policies/{id}          # Delete policy

# Policy Assignment to Devices
POST /api/v1/policies/{id}/devices    # Assign policy to device(s)
GET  /api/v1/policies/{id}/devices    # List devices with this policy
DELETE /api/v1/policies/{id}/devices/{device_id} # Remove policy from device

# Policy Assignment to Groups
POST /api/v1/policies/{id}/groups     # Assign policy to device group(s)
GET  /api/v1/policies/{id}/groups     # List groups with this policy
DELETE /api/v1/policies/{id}/groups/{group_id} # Remove policy from group
```

### 👥 Device Groups
```bash
# Group Lifecycle
POST /api/v1/device-groups            # Create device group
GET  /api/v1/device-groups            # List all groups
GET  /api/v1/device-groups/{id}       # Get group details
PUT  /api/v1/device-groups/{id}       # Update group
DELETE /api/v1/device-groups/{id}     # Delete group

# Group Membership
POST /api/v1/device-groups/{id}/devices    # Add device to group
GET  /api/v1/device-groups/{id}/devices    # List devices in group
DELETE /api/v1/device-groups/{id}/devices/{device_id} # Remove device from group
```

### 📦 Application Management
```bash
POST /api/v1/applications             # Add application
GET  /api/v1/applications             # List applications
GET  /api/v1/applications/{id}        # Get application details
PUT  /api/v1/applications/{id}        # Update application
DELETE /api/v1/applications/{id}      # Remove application
```

### 🔴 Real-time WebSocket
```bash
WS   /ws?token=JWT_TOKEN              # WebSocket connection for real-time events
GET  /api/v1/websocket/status         # WebSocket service status
```

## 🚀 API Usage Examples

### 1. Complete Device Enrollment Flow
```bash
# 1. Login to get token
TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin"}' | \
  jq -r '.token')

# 2. Enroll a new device
curl -X POST http://localhost:8081/api/v1/devices \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "uuid": "device-001",
    "hostname": "laptop-001",
    "platform": "linux",
    "os_version": "Ubuntu 22.04"
  }'

# 3. List all devices
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/devices

# 4. Execute command on device
curl -X POST http://localhost:8081/api/v1/devices/device-001/commands \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "command": "system_info",
    "parameters": {}
  }'
```

### 2. Policy Management & Assignment
```bash
# 1. Create a policy
POLICY_ID=$(curl -s -X POST http://localhost:8081/api/v1/policies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Security Policy",
    "description": "Basic security requirements",
    "type": "security",
    "rules": {
      "firewall_enabled": true,
      "auto_updates": true
    }
  }' | jq -r '.id')

# 2. Assign policy to device
curl -X POST http://localhost:8081/api/v1/policies/$POLICY_ID/devices \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"device_id": "device-001"}'

# 3. List devices with this policy
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/policies/$POLICY_ID/devices
```

### 3. Device Groups & Bulk Management
```bash
# 1. Create device group
GROUP_ID=$(curl -s -X POST http://localhost:8081/api/v1/device-groups \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Development Team",
    "description": "Developer laptops and workstations",
    "labels": {
      "department": "engineering",
      "environment": "development"
    }
  }' | jq -r '.id')

# 2. Add device to group
curl -X POST http://localhost:8081/api/v1/device-groups/$GROUP_ID/devices \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"device_id": "device-001"}'

# 3. Assign policy to entire group
curl -X POST http://localhost:8081/api/v1/policies/$POLICY_ID/groups \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"group_id": "'$GROUP_ID'"}'
```

## 🔄 Real-time Events via WebSocket

### WebSocket Connection
```javascript
// Connect to WebSocket with JWT token
const ws = new WebSocket(`ws://localhost:8081/ws?token=${jwtToken}`);

// Handle real-time events
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Real-time event:', data);
  
  switch(data.type) {
    case 'device_status_change':
      console.log(`Device ${data.data.device_id} status: ${data.data.old_status} → ${data.data.new_status}`);
      break;
    case 'policy_assignment':
      console.log(`Policy ${data.data.policy_id} ${data.data.action} to device ${data.data.device_id}`);
      break;
    case 'command_execution':
      console.log(`Command ${data.data.command} on device ${data.data.device_id}: ${data.data.status}`);
      break;
    case 'group_membership':
      console.log(`Device ${data.data.device_id} ${data.data.action} to group ${data.data.group_id}`);
      break;
  }
};
```

### Real-time Event Types
All API operations trigger corresponding real-time events:

1. **Device Operations** → `device_status_change` events
2. **Policy Assignments** → `policy_assignment` events  
3. **Command Execution** → `command_execution` events
4. **Group Changes** → `group_membership` events

## 🧪 Testing the API

### Running the Test Server
```bash
# In mobius-server directory
cd /Users/awar/Documents/Mobius/mobius-server
go run cmd/main.go serve --port 8081 --dev
```

### Running All API Tests
```bash
# Run comprehensive test suite
cd /Users/awar/Documents/Mobius/tests
./run_all_tests.sh

# Expected output: 35/35 tests passing
# - 29 Core MDM API tests
# - 6 WebSocket functionality tests
```

## 📊 API Response Formats

### Success Response
```json
{
  "id": "device-001",
  "uuid": "device-001", 
  "hostname": "laptop-001",
  "platform": "linux",
  "os_version": "Ubuntu 22.04",
  "status": "online",
  "last_seen": "2025-08-22T14:00:00Z",
  "enrolled_at": "2025-08-22T13:00:00Z",
  "labels": {}
}
```

### Error Response
```json
{
  "error": "Device not found",
  "code": "DEVICE_NOT_FOUND",
  "details": {}
}
```