#!/bin/sh
# Install peek from the rolling "latest" GitHub release.
#
#   curl -fsSL https://raw.githubusercontent.com/tatum/peek/main/install.sh | sh
#
# Installs to ~/.local/bin by default; override with PEEK_INSTALL_DIR.
set -eu

repo="tatum/peek"
dir="${PEEK_INSTALL_DIR:-$HOME/.local/bin}"

os=$(uname -s)
case "$os" in
  Linux)  os=linux ;;
  Darwin) os=darwin ;;
  *) echo "install.sh: unsupported OS: $os" >&2; exit 1 ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64|amd64)  arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo "install.sh: unsupported architecture: $arch" >&2; exit 1 ;;
esac

url="https://github.com/$repo/releases/download/latest/peek-$os-$arch"

mkdir -p "$dir"
echo "Downloading peek ($os/$arch) to $dir/peek ..."
if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$url" -o "$dir/peek"
elif command -v wget >/dev/null 2>&1; then
  wget -qO "$dir/peek" "$url"
else
  echo "install.sh: need curl or wget" >&2; exit 1
fi
chmod +x "$dir/peek"

echo "Installed: $dir/peek"
case ":$PATH:" in
  *":$dir:"*) ;;
  *) echo "Note: $dir is not on your PATH. Add it, e.g.:"
     echo "  export PATH=\"$dir:\$PATH\"" ;;
esac
