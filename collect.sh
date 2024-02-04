. ./.env;
go build cmd/huygen/main.go 

"C:\Users\Mahamed\Downloads\clumsy-0.3-win64-a\clumsy.exe" \
--filter "udp and (udp.DstPort == 8081 or udp.DstPort == 8082 or udp.DstPort == 8083)"  \
--lag on \
--lag-time 100 &

CLUMSY_PID=$!
export HUYGEN_ADDRESS=$NODE_A
export HUYGEN_PEERS=$NODE_B,$NODE_C

./main.exe &
PROC_A=$!

export HUYGEN_ADDRESS=$NODE_B
export HUYGEN_PEERS=$NODE_A,$NODE_C
./main.exe &
PROC_B=$!

export HUYGEN_ADDRESS=$NODE_C
export HUYGEN_PEERS=$NODE_B,$NODE_A
./main.exe &
PROC_C=$!

sleep 80

kill $CLUMSY_PID

sleep 1

"C:\Users\Mahamed\Downloads\clumsy-0.3-win64-a\clumsy.exe" \
--filter "udp and (udp.DstPort == 8081 or udp.DstPort == 8082 or udp.DstPort == 8083)"  \
--lag on \
--lag-time 200 &

CLUMSY_PID=$!

sleep 80

kill $CLUMSY_PID

sleep 1

"C:\Users\Mahamed\Downloads\clumsy-0.3-win64-a\clumsy.exe" \
--filter "udp and (udp.DstPort == 8081 or udp.DstPort == 8082 or udp.DstPort == 8083)"  \
--lag on \
--lag-time 300 &

CLUMSY_PID=$!

sleep 80

kill $CLUMSY_PID

sleep 1

"C:\Users\Mahamed\Downloads\clumsy-0.3-win64-a\clumsy.exe" \
--filter "udp and (udp.DstPort == 8081 or udp.DstPort == 8082 or udp.DstPort == 8083)"  \
--lag on \
--lag-time 100 &

CLUMSY_PID=$!

sleep 80


kill $CLUMSY_PID;
kill $PROC_A;
kill $PROC_B;
kill $PROC_C;
