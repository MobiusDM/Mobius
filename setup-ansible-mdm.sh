#!/bin/bash

# Mobius MDM Setup Script
# Set# Step 5: Build the Mobius applications
echo "🔧 Building Mobius applications..."

if command -v go &> /dev/null; then
    echo "Building mobius server..."
    cd mobius-server
    go build -o ../build/mobius ./cmd/mobius
    cd ..
    
    echo "Building mobiuscli..."
    cd mobius-cli
    go build -o ../build/mobiuscli ./cmd/mobiuscli
    cd ..
    
    echo "✅ Build complete"
else
    echo "❌ Please install Go: https://golang.org/"
    exit 1
fithe open-source Mobius MDM platform

set -e

echo "🚀 Setting up Mobius MDM Platform..."

# Step 1: Clean up any remaining references (this step is now complete)
echo "📝 All enterprise edition references have been removed..."

# Step 2: Install Ansible if not present
echo "🔧 Checking Ansible installation..."
if ! command -v ansible &> /dev/null; then
    echo "Installing Ansible..."
    if command -v pip3 &> /dev/null; then
        pip3 install ansible
    elif command -v brew &> /dev/null; then
        brew install ansible
    else
        echo "❌ Please install Ansible manually: https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html"
        exit 1
    fi
else
    echo "✅ Ansible is already installed"
fi

# Step 3: Set up environment
echo "⚙️  Setting up environment..."
if [ ! -f ".env" ]; then
    cp .env.example .env
    echo "📄 Created .env file. Please edit it with your configuration."
fi

# Step 4: Set up Ansible inventory
echo "📋 Setting up Ansible inventory..."
cd ansible-mdm
if [ ! -f "inventory" ]; then
    cp inventory.example inventory
    echo "📄 Created Ansible inventory. Please edit ansible-mdm/inventory with your device information."
fi
cd ..

# Step 5: Build the backend application
echo "� Building Mobius MDM backend..."
cd backend
if command -v go &> /dev/null; then
    echo "Building mobius server..."
    go build -o ../build/mobius ./cmd/mobius
    echo "Building mobiuscli..."
    go build -o ../build/mobiuscli ./cmd/mobiuscli
    echo "✅ Backend build complete"
else
    echo "❌ Please install Go: https://golang.org/"
    exit 1
fi
cd ..

# Step 6: Database setup instructions
echo "🗄️  Database setup required..."
echo "Please run the following commands to set up your database:"
echo "1. Start your database: docker compose up -d"
echo "2. Initialize Mobius database: ./build/mobius prepare db --dev"

echo ""
echo "🎉 Setup complete! Next steps:"
echo ""
echo "1. 📝 Edit .env with your configuration"
echo "2. 📋 Edit ansible-mdm/inventory with your devices (Ubuntu and Pop!_OS support included)"
echo "3. 🗄️  Set up your database (see instructions above)"
echo "4. 🚀 Start Mobius server: ./build/mobius serve --dev"
echo "5. 🌐 Access API at: http://localhost:8080/api"
echo "6. ⚙️  Run Ansible playbook: cd ansible-mdm && ansible-playbook -i inventory site.yml"
echo ""
echo "📚 For detailed documentation, see: docs/"
echo "🐛 Report issues at: https://github.com/your-repo/issues"
echo ""
echo "Happy device management with Mobius MDM! 🎯"
