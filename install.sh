#!/bin/bash
set -e

echo "🚀 Iniciando instalação dos Dotfiles..."

# 1. Instalar Starship (se não existir)
if ! command -v starship &> /dev/null; then
    echo "Instalando Starship..."
    curl -sS https://starship.rs/install.sh | sh -s -- -y
else
    echo "Starship já instalado."
fi

# 2. Configurar Starship
echo "Configurando Starship..."
mkdir -p ~/.config
# O VS Code clona o repo em ~/dotfiles por padrão
ln -sf ~/dotfiles/starship.toml ~/.config/starship.toml

# 3. Adicionar init ao .bashrc (se não estiver lá)
if ! grep -q "starship init bash" ~/.bashrc; then
    echo 'eval "$(starship init bash)"' >> ~/.bashrc
    echo "Starship adicionado ao .bashrc"
fi

# 4. Adicionar init ao .zshrc (se zsh estiver instalado)
if [ -f ~/.zshrc ]; then
    if ! grep -q "starship init zsh" ~/.zshrc; then
        echo 'eval "$(starship init zsh)"' >> ~/.zshrc
        echo "Starship adicionado ao .zshrc"
    fi
fi

echo "✅ Dotfiles instalados com sucesso!"
