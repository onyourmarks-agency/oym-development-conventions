#!/usr/bin/env bash
# Installs the latest oym-conventions release binary.
#
#   curl -fsSL https://raw.githubusercontent.com/onyourmarks-agency/oym-development-conventions/main/install.sh | bash
#
# Options via environment:
#   INSTALL_DIR   target directory (default: /usr/local/bin, falls back to ~/.local/bin)
#   GITHUB_TOKEN  used for API/download auth when the gh CLI is not installed
set -euo pipefail

REPO="onyourmarks-agency/oym-development-conventions"
BINARY="oym-conventions"

case "$(uname -s)" in
    Darwin) OS="darwin" ;;
    Linux) OS="linux" ;;
    *) echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac
case "$(uname -m)" in
    arm64|aarch64) ARCH="arm64" ;;
    x86_64|amd64) ARCH="amd64" ;;
    *) echo "unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

auth_curl() {
    if [ -n "${GITHUB_TOKEN:-}" ]; then
        curl -fsSL -H "Authorization: Bearer ${GITHUB_TOKEN}" "$@"
    else
        curl -fsSL "$@"
    fi
}

# gh refuses to run unauthenticated even against public repos, so only use it
# when it can actually authenticate.
USE_GH=false
if command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1; then
    USE_GH=true
fi

if [ "${USE_GH}" = true ]; then
    TAG="$(gh api "repos/${REPO}/releases/latest" --jq .tag_name)"
else
    TAG="$(auth_curl "https://api.github.com/repos/${REPO}/releases/latest" | grep -m1 '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
fi
if [ -z "${TAG}" ]; then
    echo "could not resolve the latest release tag" >&2
    exit 1
fi

ASSET="${BINARY}-${TAG}-${OS}-${ARCH}.tar.gz"
WORKDIR="$(mktemp -d)"
trap 'rm -rf "${WORKDIR}"' EXIT

echo "downloading ${ASSET} (${TAG})"
if [ "${USE_GH}" = true ]; then
    gh release download "${TAG}" --repo "${REPO}" --pattern "${ASSET}" --output "${WORKDIR}/${ASSET}"
elif [ -n "${GITHUB_TOKEN:-}" ]; then
    # Private repos: the browser download URL ignores tokens; go through the asset API.
    ASSET_URL="$(auth_curl "https://api.github.com/repos/${REPO}/releases/tags/${TAG}" | grep -B3 "\"name\": \"${ASSET}\"" | grep -m1 '"url"' | sed -E 's/.*"url": *"([^"]+)".*/\1/')"
    if [ -z "${ASSET_URL}" ]; then
        echo "asset ${ASSET} not found in release ${TAG}" >&2
        exit 1
    fi
    curl -fsSL -H "Authorization: Bearer ${GITHUB_TOKEN}" -H "Accept: application/octet-stream" -o "${WORKDIR}/${ASSET}" "${ASSET_URL}"
else
    auth_curl -o "${WORKDIR}/${ASSET}" "https://github.com/${REPO}/releases/download/${TAG}/${ASSET}"
fi
tar -xzf "${WORKDIR}/${ASSET}" -C "${WORKDIR}"

DEST="${INSTALL_DIR:-/usr/local/bin}"
if [ ! -w "${DEST}" ]; then
    DEST="${HOME}/.local/bin"
    mkdir -p "${DEST}"
fi
install -m 0755 "${WORKDIR}/${BINARY}" "${DEST}/${BINARY}"

echo "installed $("${DEST}/${BINARY}" version) to ${DEST}/${BINARY}"
case ":${PATH}:" in
    *":${DEST}:"*) ;;
    *) echo "note: add ${DEST} to your PATH" ;;
esac
