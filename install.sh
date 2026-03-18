#!/bin/sh
set -e

REPO="robinojw/roster"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

main() {
    detect_platform
    fetch_latest_version
    download_and_install
    configure_hooks
    echo ""
    echo "roster ${VERSION} installed to ${INSTALL_DIR}/roster"
    echo "Hooks configured for: ${CONFIGURED_HOOKS}"
    echo ""
    echo "Run 'roster bootstrap' in any repo, or start a new Claude Code session."
}

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "${ARCH}" in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *) echo "Unsupported architecture: ${ARCH}" && exit 1 ;;
    esac

    case "${OS}" in
        linux|darwin) ;;
        *) echo "Unsupported OS: ${OS}" && exit 1 ;;
    esac
}

fetch_latest_version() {
    VERSION=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "${VERSION}" ]; then
        echo "Failed to fetch latest version" && exit 1
    fi
    VERSION_NUM="${VERSION#v}"
}

download_and_install() {
    ARCHIVE="roster_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
    URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"

    echo "Downloading roster ${VERSION} for ${OS}/${ARCH}..."
    TMPDIR=$(mktemp -d)
    trap 'rm -rf "${TMPDIR}"' EXIT

    curl -sSfL "${URL}" -o "${TMPDIR}/${ARCHIVE}"
    tar -xzf "${TMPDIR}/${ARCHIVE}" -C "${TMPDIR}"

    if [ -w "${INSTALL_DIR}" ]; then
        cp "${TMPDIR}/roster" "${INSTALL_DIR}/roster"
    else
        echo "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo cp "${TMPDIR}/roster" "${INSTALL_DIR}/roster"
    fi
    chmod +x "${INSTALL_DIR}/roster"
}

configure_hooks() {
    CONFIGURED_HOOKS=""

    configure_claude_code && CONFIGURED_HOOKS="${CONFIGURED_HOOKS} claude-code"
    configure_codex && CONFIGURED_HOOKS="${CONFIGURED_HOOKS} codex"
    configure_opencode && CONFIGURED_HOOKS="${CONFIGURED_HOOKS} opencode"

    if [ -z "${CONFIGURED_HOOKS}" ]; then
        CONFIGURED_HOOKS="none (run 'roster bootstrap' manually)"
    fi
}

configure_claude_code() {
    CLAUDE_SETTINGS="${HOME}/.claude/settings.json"
    if [ ! -d "${HOME}/.claude" ]; then
        return 1
    fi

    HOOK_CMD="roster bootstrap"

    if [ ! -f "${CLAUDE_SETTINGS}" ]; then
        cat > "${CLAUDE_SETTINGS}" << 'SETTINGS'
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "startup",
        "hooks": [
          {
            "type": "command",
            "command": "roster bootstrap"
          }
        ]
      }
    ]
  }
}
SETTINGS
        return 0
    fi

    if grep -q "roster bootstrap" "${CLAUDE_SETTINGS}" 2>/dev/null; then
        return 0
    fi

    python3 -c "
import json, sys

with open('${CLAUDE_SETTINGS}', 'r') as f:
    settings = json.load(f)

hook_entry = {
    'matcher': 'startup',
    'hooks': [{'type': 'command', 'command': 'roster bootstrap'}]
}

if 'hooks' not in settings:
    settings['hooks'] = {}

if 'SessionStart' not in settings['hooks']:
    settings['hooks']['SessionStart'] = []

settings['hooks']['SessionStart'].append(hook_entry)

with open('${CLAUDE_SETTINGS}', 'w') as f:
    json.dump(settings, f, indent=2)
    f.write('\n')
" 2>/dev/null || {
        echo "Warning: Could not auto-configure Claude Code hook."
        echo "Add this to ${CLAUDE_SETTINGS} manually:"
        echo '  "hooks": { "SessionStart": [{ "matcher": "startup", "hooks": [{ "type": "command", "command": "roster bootstrap" }] }] }'
        return 1
    }
    return 0
}

configure_codex() {
    CODEX_CONFIG="${HOME}/.codex"
    if [ ! -d "${CODEX_CONFIG}" ]; then
        return 1
    fi

    CODEX_INSTRUCTIONS="${CODEX_CONFIG}/instructions.md"
    if [ -f "${CODEX_INSTRUCTIONS}" ] && grep -q "roster bootstrap" "${CODEX_INSTRUCTIONS}" 2>/dev/null; then
        return 0
    fi

    cat >> "${CODEX_INSTRUCTIONS}" << 'CODEX'

## Roster

Run `roster bootstrap` at the start of each session to load agent personas and repo signals. The output is written to `AGENTS.md` and `.roster/`.
CODEX
    return 0
}

configure_opencode() {
    OPENCODE_CONFIG="${HOME}/.config/opencode"
    if [ ! -d "${OPENCODE_CONFIG}" ]; then
        return 1
    fi

    OPENCODE_INSTRUCTIONS="${OPENCODE_CONFIG}/instructions.md"
    if [ -f "${OPENCODE_INSTRUCTIONS}" ] && grep -q "roster bootstrap" "${OPENCODE_INSTRUCTIONS}" 2>/dev/null; then
        return 0
    fi

    cat >> "${OPENCODE_INSTRUCTIONS}" << 'OPENCODE'

## Roster

Run `roster bootstrap` at the start of each session to load agent personas and repo signals. The output is written to `AGENTS.md` and `.roster/`.
OPENCODE
    return 0
}

main
