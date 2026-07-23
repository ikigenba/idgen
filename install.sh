#!/bin/sh

set -eu

REPO="ikigenba/idgen"
BINARY="idgen"
VERSION="${IDGEN_VERSION:-latest}"
DESTINATION="${BINDIR:-${PREFIX:-$HOME/.local}/bin}"

case "$(uname -s)" in
	Linux) os="linux" ;;
	Darwin) os="darwin" ;;
	*)
		echo "idgen: unsupported operating system: $(uname -s)" >&2
		exit 1
		;;
esac

case "$(uname -m)" in
	x86_64 | amd64) arch="amd64" ;;
	arm64 | aarch64) arch="arm64" ;;
	*)
		echo "idgen: unsupported architecture: $(uname -m)" >&2
		exit 1
		;;
esac

if [ "$VERSION" = "latest" ]; then
	release_path="releases/latest/download"
else
	release_path="releases/download/$VERSION"
fi

archive="${BINARY}_${os}_${arch}.tar.gz"
url="https://github.com/${REPO}/${release_path}/${archive}"
temporary="$(mktemp -d)"
trap 'rm -rf "$temporary"' EXIT HUP INT TERM

curl -fsSL "$url" -o "$temporary/$archive"
tar -xzf "$temporary/$archive" -C "$temporary"
install -d "$DESTINATION"
install -m 0755 "$temporary/$BINARY" "$DESTINATION/$BINARY"

echo "Installed $BINARY to $DESTINATION/$BINARY"
