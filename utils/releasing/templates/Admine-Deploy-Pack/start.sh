cd pubsub/redis
docker compose up -d
sleep 2
cd ../../server_handler
./server_handler &
cd ../vpn_handler
./vpn_handler &
cd ../bot
./bot &
