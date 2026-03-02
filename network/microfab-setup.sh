
#!/bin/bash
set -e

# Diretório raiz do repositório (pai do diretório network/)
REPO_DIR="$(cd "$(dirname "$0")/.." && pwd)"
NETWORK_DIR="$REPO_DIR/network"
WALLET_DIR="$NETWORK_DIR/wallet/org1admin"
MSP_DIR="$WALLET_DIR/msp"
BIN_DIR="$REPO_DIR/bin"
MICROFAB_URL="http://localhost:8080"

echo "==> [1/5] Subindo o container Microfab..."
if [ "$(docker ps -q -f name=microfab)" ]; then
  echo "    Container já está rodando, pulando..."
else
  docker run -d \
    --name microfab \
    -p 8080:8080 \
    -e MICROFAB_CONFIG='{
      "endorsing_organizations": [{ "name": "Org1" }],
      "channels": [{ "name": "mychannel", "endorsing_organizations": ["Org1"] }],
      "capability_level": "V2_0"
    }' \
    ibmcom/ibp-microfab:latest

  echo "    Aguardando Microfab inicializar..."
  sleep 5
fi

echo "==> [2/5] Baixando binários do Hyperledger Fabric..."
if [ -f "$BIN_DIR/peer" ]; then
  echo "    Binários já existem, pulando..."
else
  cd "$REPO_DIR"
  curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh \
    | bash -s -- binary
fi

echo "==> [3/5] Extraindo identidades do Microfab..."
mkdir -p "$MSP_DIR"/{signcerts,keystore,cacerts,admincerts}

curl -s "$MICROFAB_URL/ak/api/v1/components" | \
  jq -r '.[] | select(.id == "org1admin") | .cert' | base64 -d \
  > "$MSP_DIR/signcerts/cert.pem"

curl -s "$MICROFAB_URL/ak/api/v1/components" | \
  jq -r '.[] | select(.id == "org1admin") | .private_key' | base64 -d \
  > "$MSP_DIR/keystore/key.pem"

curl -s "$MICROFAB_URL/ak/api/v1/components" | \
  jq -r '.[] | select(.id == "org1admin") | .ca' | base64 -d \
  > "$MSP_DIR/cacerts/ca.pem"

cp "$MSP_DIR/signcerts/cert.pem" "$MSP_DIR/admincerts/cert.pem"

echo "==> [4/5] Criando config.yaml do MSP..."
cat > "$MSP_DIR/config.yaml" << 'EOF'
NodeOUs:
  Enable: false
EOF

echo "==> [5/5] Gerando fabric-env.sh..."
cat > "$NETWORK_DIR/fabric-env.sh" << EOF
export PATH=$BIN_DIR:\$PATH
export FABRIC_CFG_PATH=$REPO_DIR/config/
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=$MSP_DIR
export CORE_PEER_ADDRESS=org1peer-api.127-0-0-1.nip.io:8080
export CORE_PEER_TLS_ENABLED=false
EOF

echo ""
echo "✅ Setup concluído!"
echo ""
echo "Para ativar o ambiente, execute:"
echo "  source network/fabric-env.sh"
echo ""
echo "Para testar, execute:"
echo "  peer channel list"

