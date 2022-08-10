# Cross compile to raspberry pi
env GOOS=linux GOARCH=arm GOARM=5 go build -o ./dist/pigeon
# Transfer pigeon service executable and unit file securely
scp ./dist/pigeon ./dist/pigeon.service admin@pi:~/pigeon/
# Make service executable, move unit file to /lib/systemd/system/ and set permissions
ssh admin@pi '\
chmod +x ~/pigeon/pigeon; \
sudo mv ~/pigeon/pigeon.service /lib/systemd/system/pigeon.service; \
sudo chmod 740 /lib/systemd/system/pigeon.service'