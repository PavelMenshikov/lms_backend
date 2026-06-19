#!/usr/bin/env bash
set -euo pipefail

REPO_URL="https://github.com/GrtsqDev/cap-ed-backend.git"
PROJECT_DIR="cap-ed-backend"

GOLANG_VERSION="1.25.0"

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

log()  { echo -e "${CYAN}[INFO]${NC} $*"; }
ok()   { echo -e "${GREEN}[OK]${NC}   $*"; }
err()  { echo -e "${RED}[ERR]${NC}  $*"; }

detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        OS_LIKE=$ID_LIKE
    elif command -v lsb_release &>/dev/null; then
        OS=$(lsb_release -si | tr '[:upper:]' '[:lower:]')
    else
        OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    fi
}

install_system_packages() {
    log "Installing system packages..."
    case "$OS" in
        ubuntu|debian|linuxmint|pop)
            sudo apt-get update -qq
            sudo apt-get install -y -qq curl wget git ca-certificates gnupg lsb-release
            ;;
        fedora|rhel|centos)
            sudo yum install -y -q curl wget git ca-certificates
            ;;
        arch|manjaro)
            sudo pacman -Sy --noconfirm curl wget git ca-certificates
            ;;
        alpine)
            sudo apk add --no-cache curl wget git ca-certificates
            ;;
        *)
            err "Unsupported distro: $OS. Install curl, wget, git manually."
            exit 1
            ;;
    esac
    ok "System packages installed"
}

install_docker() {
    if command -v docker &>/dev/null; then
        log "Docker already installed: $(docker --version)"
        return
    fi
    log "Installing Docker..."
    case "$OS" in
        ubuntu|debian|linuxmint|pop)
            curl -fsSL https://get.docker.com | sudo bash
            ;;
        fedora|rhel|centos)
            curl -fsSL https://get.docker.com | sudo bash
            ;;
        arch|manjaro)
            sudo pacman -Sy --noconfirm docker docker-compose
            sudo systemctl enable --now docker
            ;;
        alpine)
            sudo apk add --no-cache docker docker-compose
            sudo rc-update add docker default
            sudo service docker start
            ;;
        *)
            curl -fsSL https://get.docker.com | sudo bash
            ;;
    esac
    sudo usermod -aG docker "$USER"
    ok "Docker installed. You may need to re-login for group changes."
}

install_docker_compose() {
    if docker compose version &>/dev/null 2>&1; then
        log "Docker Compose already installed: $(docker compose version)"
        return
    fi
    if command -v docker-compose &>/dev/null; then
        log "docker-compose (v1) already installed"
        return
    fi
    log "Installing Docker Compose plugin..."
    case "$OS" in
        ubuntu|debian|linuxmint|pop)
            sudo apt-get install -y -qq docker-compose-plugin 2>/dev/null || {
                DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
                mkdir -p "$DOCKER_CONFIG/cli-plugins"
                sudo curl -SL "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
                    -o "$DOCKER_CONFIG/cli-plugins/docker-compose"
                sudo chmod +x "$DOCKER_CONFIG/cli-plugins/docker-compose"
            }
            ;;
        *)
            DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
            mkdir -p "$DOCKER_CONFIG/cli-plugins"
            sudo curl -SL "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
                -o "$DOCKER_CONFIG/cli-plugins/docker-compose"
            sudo chmod +x "$DOCKER_CONFIG/cli-plugins/docker-compose"
            ;;
    esac
    ok "Docker Compose installed"
}

install_golang() {
    if command -v go &>/dev/null; then
        local current_version
        current_version=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
        if [ "$(echo "$current_version >= 1.25" | bc -l 2>/dev/null || echo 0)" = 1 ]; then
            log "Go already installed: $(go version)"
            return
        fi
        log "Upgrading Go from $current_version to $GOLANG_VERSION..."
    fi

    log "Installing Go $GOLANG_VERSION..."
    local arch
    arch=$(uname -m)
    case "$arch" in
        x86_64)  arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        *)       err "Unsupported architecture: $arch"; exit 1 ;;
    esac

    local tarball="go${GOLANG_VERSION}.linux-${arch}.tar.gz"
    wget -q "https://go.dev/dl/${tarball}" -O "/tmp/${tarball}"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "/tmp/${tarball}"
    rm -f "/tmp/${tarball}"

    if ! grep -q '/usr/local/go/bin' "$HOME/.profile" 2>/dev/null; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> "$HOME/.profile"
    fi
    export PATH=$PATH:/usr/local/go/bin
    ok "Go $GOLANG_VERSION installed"
}

clone_repo() {
    if [ -d "$PROJECT_DIR" ]; then
        log "Project directory '$PROJECT_DIR' already exists, pulling..."
        cd "$PROJECT_DIR"
        git pull
        cd ..
    else
        log "Cloning repository..."
        git clone "$REPO_URL"
        ok "Repository cloned"
    fi
}

setup_env() {
    cd "$PROJECT_DIR"
    if [ ! -f .env ]; then
        if [ -f .env.example ]; then
            cp .env.example .env
            log "Created .env from .env.example"
        else
            cat > .env << 'ENVEOF'
DB_HOST=db
DB_PORT=5432
DB_USER=capuser
DB_PASSWORD=cappass
DB_NAME=capedu
API_PORT=8000
REDIS_ADDR=redis:6379
S3_ENDPOINT_URL=https://fs.cloupard.kz
S3_REGION=us-east-1
S3_BUCKET_NAME=capeducation
S3_ACCESS_KEY_ID=JLC3NTLH51VE0LB7D1EZ
S3_SECRET_ACCESS_KEY=qu3BqIk4D90Dhz9piWIDSwBoRrCukhA84HHATwB1
SYSTEM_SECRET=my_ultra_secret_backdoor_key_2026
ENVEOF
            log "Created .env with default values"
        fi
    else
        log ".env already exists, keeping it"
    fi
    cd ..
}

run_docker() {
    cd "$PROJECT_DIR"
    log "Starting services with Docker Compose..."
    sudo docker compose up -d
    ok "All services started!"
    echo ""
    echo "============================================"
    echo "  Backend API:  http://localhost:8000"
    echo "  Swagger docs: http://localhost:8000/swagger/"
    echo "============================================"
    echo ""
    log "Useful commands:"
    echo "  docker compose logs -f       # follow logs"
    echo "  docker compose down          # stop services"
    echo "  docker compose restart app   # restart backend only"
    echo ""
}

print_help() {
    cat << HELP
Usage: $0 [OPTIONS]

Setup script for Cap Education LMS Backend.

Options:
  --skip-docker    Skip Docker installation
  --skip-go        Skip Go installation
  --skip-clone     Skip repository clone
  --no-run         Don't start services after setup
  -h, --help       Show this help
HELP
    exit 0
}

SKIP_DOCKER=false
SKIP_GO=false
SKIP_CLONE=false
NO_RUN=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        --skip-docker) SKIP_DOCKER=true; shift ;;
        --skip-go)     SKIP_GO=true; shift ;;
        --skip-clone)  SKIP_CLONE=true; shift ;;
        --no-run)      NO_RUN=true; shift ;;
        -h|--help)     print_help ;;
        *)             echo "Unknown option: $1"; print_help ;;
    esac
done

echo ""
echo "============================================"
echo "  Cap Education — Client Setup"
echo "============================================"
echo ""

detect_distro
install_system_packages

if [ "$SKIP_DOCKER" = false ]; then
    install_docker
    install_docker_compose
fi

if [ "$SKIP_GO" = false ]; then
    install_golang
fi

if [ "$SKIP_CLONE" = false ]; then
    clone_repo
fi

setup_env

if [ "$NO_RUN" = false ]; then
    run_docker
fi

echo ""
ok "Setup complete!"
