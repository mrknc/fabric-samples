from fabric-samples:
cd basic-network
./start.sh
cd ../commercial-paper/organization/magnetocorp/configuration/cli
docker-compose -f docker-compose.yml up -d cliMagnetoCorp
./monitordocker.sh net_basic

docker exec cliMagnetoCorp peer chaincode install -n papercontract -v 0 -p github.com/contract-go -l golang
docker exec cliMagnetoCorp peer chaincode instantiate -n papercontract -v 0 -l golang -c '{"Args":["org.papernet.commercialpaper:instantiate"]}' -C mychannel -P "AND ('Org1MSP.member')"
docker exec cliMagnetoCorp peer chaincode invoke -n papercontract -c '{"Args":["issuePaper", "MagnetoCorp","00001", "2020-01-01","2020-07-01","1000"]}' -C mychannel
docker exec cliMagnetoCorp peer chaincode invoke -n papercontract -c '{"Args":["buyPaper", "MagnetoCorp","00001", "MagnetoCorp","DigiBank"]}' -C mychannel
docker exec cliMagnetoCorp peer chaincode invoke -n papercontract -c '{"Args":["redeemPaper", "MagnetoCorp","00001","DigiBank"]}' -C mychannel

DB browser: http://127.0.0.1:5984/_utils/#database/mychannel_papercontract/_all_docs