#!/bin/bash

# Function to check if the script is being run as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo
        echo "Please run this script as root."
        echo
        sleep 1
        exit 1
    fi
}

# Function to install dependencies based on package manager
install_package() {
    # Detect package manager
    if command -v apt > /dev/null; then
        PKG_MANAGER="apt"
    elif command -v yum > /dev/null; then
        PKG_MANAGER="yum"
    else
        echo "Unsupported package manager. Exiting."
        exit 1
    fi

    echo
    echo "Updating the OS and installing dependencies..."
    echo
    sleep 0.5

    if [ "$PKG_MANAGER" == "apt" ]; then
        sudo apt update -q
        sudo apt upgrade -y
        sudo apt install -y curl wget build-essential ca-certificates git
    elif [ "$PKG_MANAGER" == "yum" ]; then
        sudo yum update -y
        sudo yum groupinstall -y "Development Tools"
        sudo yum install -y curl wget git ca-certificates
    fi

    echo
    echo "OS Updated & Dependencies Installed."
    echo
    sleep 0.5
}

# Function to download and install Go
install_go() {
    local go_version="go1.21.0"
    echo
    echo "Installing Go $go_version..."
    echo
    sleep 0.5
    local go_url="https://go.dev/dl/$go_version.linux-amd64.tar.gz"
    sleep 0.5

    # Remove any existing Go installation
    sudo rm -rf /usr/local/go
    sleep 0.5

    # Download and extract the Go archive
    sudo curl -sLo go.tar.gz "$go_url"
    sleep 0.5

    sudo tar -C /usr/local/ -xzf go.tar.gz
    sleep 0.5

    sudo rm go.tar.gz

    # Add Go binary path to system PATH
    export PATH=$PATH:/usr/local/go/bin
    echo "export PATH=\$PATH:/usr/local/go/bin" > /etc/profile.d/go.sh
    sleep 0.5

    source /etc/profile.d/go.sh
    sleep 0.5

    # Check installed Go version
    go version
    sleep 0.5

    echo
    echo "Go $go_version installed successfully."
    echo
}

# Check if running as root
check_root

# Install dependencies for Go and Git
install_package

# Install Go
install_go
