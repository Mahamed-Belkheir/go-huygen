go build ./cmd/sensor/
go build ./cmd/follower/

NODE_A_ADDRESS_SENSOR=127.0.0.1:3001
NODE_A_ADDRESS_FOLLOWER=127.0.0.1:3011
NODE_A_PEER_LATENCY=127.0.0.1:3002=100,127.0.0.1:3003=200

NODE_B_ADDRESS_SENSOR=127.0.0.1:3002
NODE_B_ADDRESS_FOLLOWER=127.0.0.1:3012
NODE_B_PEER_LATENCY=127.0.0.1:3001=100,127.0.0.1:3003=200

NODE_C_ADDRESS_SENSOR=127.0.0.1:3003
NODE_C_ADDRESS_FOLLOWER=127.0.0.1:3013
NODE_C_PEER_LATENCY=127.0.0.1:3001=200,127.0.0.1:3002=200

export HUYGENS_PEERS_LATENCY=$NODE_A_PEER_LATENCY
export HUYGEN_ADDRESS=$NODE_A_ADDRESS_SENSOR
echo "running sensor A"
./sensor.exe &
A_SENSOR=$!


export HUYGEN_ADDRESS=$NODE_A_ADDRESS_FOLLOWER
export HUYGEN_SENSOR_ADDRESS=$NODE_A_ADDRESS_SENSOR
export HUYGEN_PEERS=$NODE_B_ADDRESS_SENSOR,$NODE_C_ADDRESS_SENSOR
export HUYGENS_PEERS_LATENCY=""
echo "running follower A"
./follower.exe &
A_FOLLOWER=$!

export HUYGENS_PEERS_LATENCY=$NODE_B_PEER_LATENCY
export HUYGEN_ADDRESS=$NODE_B_ADDRESS_SENSOR
echo "running sensor B"
./sensor.exe &
B_SENSOR=$!


export HUYGEN_ADDRESS=$NODE_B_ADDRESS_FOLLOWER
export HUYGEN_SENSOR_ADDRESS=$NODE_B_ADDRESS_SENSOR
export HUYGEN_PEERS=$NODE_A_ADDRESS_SENSOR,$NODE_C_ADDRESS_SENSOR
export HUYGENS_PEERS_LATENCY=""
echo "running follower B"
./follower.exe &
B_FOLLOWER=$!


export HUYGENS_PEERS_LATENCY=$NODE_C_PEER_LATENCY
export HUYGEN_ADDRESS=$NODE_C_ADDRESS_SENSOR
echo "running sensor C"
./sensor.exe &
C_SENSOR=$!


export HUYGEN_ADDRESS=$NODE_C_ADDRESS_FOLLOWER
export HUYGEN_SENSOR_ADDRESS=$NODE_C_ADDRESS_SENSOR
export HUYGEN_PEERS=$NODE_A_ADDRESS_SENSOR,$NODE_B_ADDRESS_SENSOR
export HUYGENS_PEERS_LATENCY=""
echo "running follower C"
./follower.exe &
C_FOLLOWER=$!

read var