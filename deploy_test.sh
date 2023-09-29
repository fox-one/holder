set -e

env GOOS=linux GOARCH=amd64 task build
scp ./builds/holder-worker ptest-dolphin:~/hodl/worker.new
ssh ptest-dolphin "cd ~/hodl; sudo systemctl stop hodl; cp worker.new worker; sudo systemctl restart hodl"
