#!/usr/bin/env bash
set -euo pipefail

# GitHub Org and Repo to get archives from
GITHUB_ORG="micro"
GITHUB_REPO="micro"

# micro install directory
MICRO_INSTALL_DIR="/usr/local/bin"
# micro cli name
MICRO_CLI_NAME="micro"
# micro cli install path
MICRO_CLI_PATH="${MICRO_INSTALL_DIR}/${MICRO_CLI_NAME}"

# get machine ARCH
ARCH=$(uname -m)
# get machine OS
OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')
# Linux requires sudo for $MICRO_INSTALL_DIR
SUDO="false"

# Http request CLI
HTTP_CLIENT=curl

getSystemInfo() {
    echo "Getting system information"
    case $ARCH in
        armv7*)
                ARCH="arm";;
        aarch64)
                ARCH="arm64";;
        x86_64)
                ARCH="amd64";;
    esac

    # linux requires sudo permissions
    if [ "$OS" == "linux" ]; then
        SUDO="true"
    fi
    echo "Your machine is running ${OS} on ${ARCH} CPU architecture"
}

checkSupported() {
    local supported_osarch=(darwin-amd64 linux-amd64 linux-arm7 linux-arm64)
    local machine_osarch="${OS}-${ARCH}"

    echo "Checking machine system support"
    for osarch in "${supported_osarch[@]}"; do
        if [ "$osarch" == "$machine_osarch" ]; then
            return
        fi
    done

    echo "No prebuilt binary for ${machine_osarch}"
    exit 1
}

checkHttpClient() {
    echo "Checking HTTP client"
    if type "curl" > /dev/null; then
        HTTP_CLIENT="curl"
    elif type "wget" > /dev/null; then
        HTTP_CLIENT="wget"
    else
        echo "Either curl or wget is required"
        exit 1
    fi
}

sudoRun() {
    local CMD="$*"

    if [ $EUID -ne 0 -a $SUDO = "true" ]; then
        CMD="sudo $CMD"
    fi

    $CMD
}

getLatestRelease() {
    local release_url="https://api.github.com/repos/${GITHUB_ORG}/${GITHUB_REPO}/releases"
    local latest_release=""

    echo "Getting the latest micro release"
    if [ "$HTTP_CLIENT" == "curl" ]; then
        latest_release=$(curl -s $release_url | grep \"tag_name\" | awk 'NR==1{print $2}' |  sed -n 's/\"\(.*\)\",/\1/p')
    else
        latest_release=$(wget -q --header="Accept: application/json" -O - $release_url | grep \"tag_name\" | awk 'NR==1{print $2}' |  sed -n 's/\"\(.*\)\",/\1/p')
    fi
    echo "Latest micro release found: ${latest_release}"

    LATEST_RELEASE_TAG=$latest_release
    CLI_ARCHIVE="${MICRO_CLI_NAME}-${LATEST_RELEASE_TAG}-${OS}-${ARCH}.tar.gz"
    DOWNLOAD_BASE="https://github.com/${GITHUB_ORG}/${GITHUB_REPO}/releases/download"
    DOWNLOAD_URL="${DOWNLOAD_BASE}/${LATEST_RELEASE_TAG}/${CLI_ARCHIVE}"

    TMP_ROOT=$(mktemp -dt micro-install-XXXXXX)
    TMP_FILE="$TMP_ROOT/$CLI_ARCHIVE"

    echo "Downloading $DOWNLOAD_URL ..."
    if [ "$HTTP_CLIENT" == "curl" ]; then
        curl -SsL "$DOWNLOAD_URL" -o "$TMP_FILE"
    else
        wget -q -O "$TMP_FILE" "$DOWNLOAD_URL"
    fi

    if [ ! -f "$TMP_FILE" ]; then
        echo "Failed to download $DOWNLOAD_URL ..."
        exit 1
    fi
}

installFile() {
    tar xf "$TMP_FILE" -C "$TMP_ROOT"
    local tmp_root_cli="$TMP_ROOT/$MICRO_CLI_NAME"

    if [ ! -f "$tmp_root_cli" ]; then
        echo "Failed to unpack micro cli binary."
        exit 1
    fi

    chmod o+x $tmp_root_cli
    sudoRun cp "$tmp_root_cli" "$MICRO_INSTALL_DIR"

    if [ -f "$MICRO_CLI_PATH" ]; then
        echo "$MICRO_CLI_NAME installed into $MICRO_INSTALL_DIR successfully."

        $MICRO_CLI_PATH --version
    else
        echo "Failed to install $MICRO_CLI_NAME"
        exit 1
    fi
}

fail_trap() {
    result=$?
    if [ "$result" != "0" ]; then
        echo "Failed to install micro"
        echo "For support, please file an issue in https://github.com/micro/micro/issues"
    fi
    cleanup
    exit $result
}

cleanup() {
    if [[ -d "${TMP_ROOT:-}" ]]; then
        rm -rf "$TMP_ROOT"
    fi
}

printInfo() {
    echo -e "\nTo get started with micro please visit official documentation https://micro.mu/docs"
    echo "To start contributing to micro please visit https://github.com/micro"
    echo "Join micro community on slack https://micro.mu/slack"
}

# catch errors and print help
trap "fail_trap" EXIT

# execute installation
getSystemInfo
checkSupported
checkHttpClient
getLatestRelease
installFile
cleanup
printInfo
