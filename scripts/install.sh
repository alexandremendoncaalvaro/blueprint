#!/usr/bin/env bash
set -e

# Bootstrap: clona o repo e compila o binario do dotfiles manager.
# Uso: curl -sL <url>/install.sh | bash

REPO="https://github.com/ale/dotfiles.git"
INSTALL_DIR="$HOME/dotfiles"

echo "=== Dotfiles Bootstrap ==="
echo

# 1. Verificar/instalar Go
if ! command -v go &>/dev/null; then
    echo "Go nao encontrado. Instalando..."
    curl -sL https://go.dev/dl/go1.23.6.linux-amd64.tar.gz -o /tmp/go.tar.gz
    sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    export PATH="/usr/local/go/bin:$PATH"
    echo "Go instalado: $(go version)"
fi

# 2. Clonar ou atualizar repo
if [ -d "$INSTALL_DIR/.git" ]; then
    echo "Atualizando repositorio..."
    git -C "$INSTALL_DIR" pull --ff-only
else
    echo "Clonando repositorio..."
    git clone "$REPO" "$INSTALL_DIR"
fi

# 3. Compilar
echo "Compilando..."
cd "$INSTALL_DIR"
make build

# 4. Instalar no PATH
mkdir -p "$HOME/.local/bin"
ln -sf "$INSTALL_DIR/bin/dotfiles" "$HOME/.local/bin/dotfiles"

echo
echo "Instalado! Use: dotfiles apply"
echo
echo "Garanta que ~/.local/bin esta no seu PATH:"
echo '  export PATH="$HOME/.local/bin:$PATH"'
