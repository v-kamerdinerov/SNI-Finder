#!/bin/bash

FontColor_Red="\033[31m"
FontColor_Green="\033[32m"
FontColor_Yellow="\033[33m"
FontColor_Suffix="\033[0m"


log() {
    local LEVEL="$1"
    local MSG="$2"
    case "${LEVEL}" in
    INFO)
        local LEVEL="[${FontColor_Green}${LEVEL}${FontColor_Suffix}]"
        local MSG="${LEVEL} ${MSG}"
        ;;
    WARN)
        local LEVEL="[${FontColor_Yellow}${LEVEL}${FontColor_Suffix}]"
        local MSG="${LEVEL} ${MSG}"
        ;;
    ERROR)
        local LEVEL="[${FontColor_Red}${LEVEL}${FontColor_Suffix}]"
        local MSG="${LEVEL} ${MSG}"
        ;;
    *) ;;
    esac
    echo -e "${MSG}"
}


# Function to check if the script is being run as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo
        log ERROR "Please run this script as root."
        echo
        sleep 1
        exit 1
    fi
}

# Function to check if golang is installed
check_go() {
    if command -v go > /dev/null; then
        log INFO "Golang installed!"
        go version
        sleep 1
        exit 0
    else
        log WARN "Golang not found. Proceed..."
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
        log ERROR "Unsupported package manager. Exiting."
        exit 1
    fi

    echo
    log WARN "Updating the OS and installing dependencies..."
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
    log WARN "OS Updated & Dependencies Installed."
    echo
    sleep 0.5
}

# Function to download and install Go
install_go() {
    local go_version="go1.21.0"
    echo
    log WARN "Installing Go $go_version..."
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
    log INFO "Go $go_version installed successfully."
    echo
}


# Check if running as root
check_root

# Check if golang installed
check_go

# Install dependencies for Go and Git
install_package

# Install Go
install_go
