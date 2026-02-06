#!/usr/bin/env bash
set -e

# ============================================
# Fix: ' + c => ç on Bluefin (Wayland/GNOME)
# Method: Compose (~/.XCompose) with validation
# ============================================

GREEN="\033[32m"
YELLOW="\033[33m"
RED="\033[31m"
RESET="\033[0m"

ok()   { echo -e "${GREEN}[OK]${RESET} $1"; }
warn() { echo -e "${YELLOW}[WARN]${RESET} $1"; }
fail() { echo -e "${RED}[FAIL]${RESET} $1"; exit 1; }

section() {
  echo
  echo "----------------------------------------"
  echo "$1"
  echo "----------------------------------------"
}

# ----------------------------
# 1. Check session
# ----------------------------
section "Session check"

if [ "$XDG_SESSION_TYPE" = "wayland" ]; then
  ok "Wayland session detected"
else
  warn "Session is not Wayland (X11 or unknown)"
fi


# ----------------------------
# 2. Check IBus
# ----------------------------
section "Input method check"

if pgrep -f ibus-daemon >/dev/null; then
  ok "IBus is running"
else
  warn "IBus not detected (Compose may fail in some apps)"
fi


# ----------------------------
# 3. Check keyboard layout
# ----------------------------
section "Keyboard layout check"

LAYOUT="unknown"
VARIANT="unknown"

if command -v localectl >/dev/null; then
  LAYOUT=$(localectl status 2>/dev/null | grep "X11 Layout" | awk -F: '{print $2}' | xargs || true)
  VARIANT=$(localectl status 2>/dev/null | grep "X11 Variant" | awk -F: '{print $2}' | xargs || true)
fi

echo "Layout : $LAYOUT"
echo "Variant: $VARIANT"

# Accept multiple layouts (comma separated)
if [[ "$LAYOUT" != *"us"* ]]; then
  fail "No US layout detected. Set US International first."
fi

if [[ "$VARIANT" != *"intl"* ]]; then
  fail "No intl (dead keys) variant detected."
fi


ok "US International with dead keys detected"


# ----------------------------
# 4. Prepare ~/.XCompose
# ----------------------------
section "Configuring Compose"

XCOMPOSE="$HOME/.XCompose"

# Backup
if [ -f "$XCOMPOSE" ]; then
  BACKUP="$XCOMPOSE.bak.$(date +%Y%m%d-%H%M%S)"
  cp "$XCOMPOSE" "$BACKUP"
  ok "Backup created: $BACKUP"
fi


# Ensure base file
if [ ! -f "$XCOMPOSE" ]; then
  echo 'include "%L"' > "$XCOMPOSE"
  ok "Created ~/.XCompose"
fi


# Ensure include "%L"
if ! grep -q 'include "%L"' "$XCOMPOSE"; then
  TMP=$(mktemp)

  echo 'include "%L"' > "$TMP"
  echo >> "$TMP"
  cat "$XCOMPOSE" >> "$TMP"

  mv "$TMP" "$XCOMPOSE"

  ok 'Added include "%L"'
fi


# Remove old block (if exists)
sed -i '/# BEGIN BLUEFIN CEDILLA/,/# END BLUEFIN CEDILLA/d' "$XCOMPOSE"


# Append new rules
cat >> "$XCOMPOSE" << 'EOF'

# BEGIN BLUEFIN CEDILLA
<dead_acute> <c> : "ç"
<dead_acute> <C> : "Ç"
# END BLUEFIN CEDILLA
EOF


ok "Compose rules written"


# ----------------------------
# 5. Finish
# ----------------------------
section "Done"

echo
echo "Configuration applied successfully."
echo
echo "IMPORTANT:"
echo "You must log out and log in again (or reboot)."
echo
echo "After that, test:"
echo "  ' + c  => ç"
echo "  ' + C  => Ç"
echo
