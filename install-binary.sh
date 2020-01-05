
#!/bin/sh -e

# Copied w/ love from the excellent hypnoglow/helm-s3

if [ -n "${HELM_PUSH_PLUGIN_NO_INSTALL_HOOK}" ]; then
    echo "Development mode: not downloading versioned release."
    exit 0
fi

version="$(curl -s https://api.github.com/repos/maorfr/helm-backup/releases/latest | awk -F '"' '/tag_name/ {print $4}')"
echo "Downloading and installing helm-backup ${version} ..."

url=""
if [ "$(uname)" = "Darwin" ]; then
    url="https://github.com/maorfr/helm-backup/releases/download/${version}/helm-backup-macos-${version}.tgz"
elif [ "$(uname)" = "Linux" ] ; then
    url="https://github.com/maorfr/helm-backup/releases/download/${version}/helm-backup-linux-${version}.tgz"
else
    url="https://github.com/maorfr/helm-backup/releases/download/${version}/helm-backup-windows-${version}.tgz"
fi

echo $url

cd $HELM_PLUGIN_DIR
mkdir -p "bin"
mkdir -p "releases/${version}"

# Download with curl if possible.
if [ -x "$(which curl 2>/dev/null)" ]; then
    curl -sSL "${url}" -o "releases/${version}.tgz"
else
    wget -q "${url}" -O "releases/${version}.tgz"
fi
tar xzf "releases/${version}.tgz" -C "releases/${version}"
mv "releases/${version}/backup" "bin/backup" || \
    mv "releases/${version}/backup.exe" "bin/backup"
rm -rf releases