# Battle Game

A simple real-time battle game where two players can fight each other by clicking/tapping to deal damage.

## Features
- Real-time WebSocket communication
- Random matchmaking
- Health bars and damage indicators
- Healing mechanics
- Mobile-friendly interface

## Local Development

1. Install Go (1.21 or later)
2. Clone the repository
```bash
git clone https://github.com/yourusername/golang-battle.git
cd golang-battle
```
3. Install dependencies
```bash
go mod download
```
4. Run the server
```bash
go run main.go
```
5. Open `http://localhost:8080` in your browser

## Deployment to Digital Ocean

### Method 1: Manual Deployment

1. Create a new Ubuntu droplet on Digital Ocean
2. SSH into your droplet:
```bash
ssh root@your_droplet_ip
```

3. Install Go:
```bash
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

4. Clone the repository:
```bash
git clone https://github.com/yourusername/golang-battle.git
cd golang-battle
```

5. Build and run:
```bash
go build -o battle-game
./battle-game
```

6. (Optional) Create a systemd service to run the game in the background:
```bash
sudo nano /etc/systemd/system/battle-game.service
```

Add the following content:
```ini
[Unit]
Description=Battle Game Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/golang-battle
ExecStart=/root/golang-battle/battle-game
Restart=always

[Install]
WantedBy=multi-user.target
```

Start the service:
```bash
sudo systemctl enable battle-game
sudo systemctl start battle-game
```

### Method 2: Docker Deployment

1. Install Docker on your droplet:
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
```

2. Clone the repository:
```bash
git clone https://github.com/yourusername/golang-battle.git
cd golang-battle
```

3. Build and run with Docker:
```bash
docker build -t battle-game .
docker run -d -p 80:8080 battle-game
```

### Method 3: GitHub Actions + Docker Hub (Recommended)

1. Fork this repository on GitHub

2. Add the following secrets to your GitHub repository:
   - `DOCKERHUB_USERNAME`
   - `DOCKERHUB_TOKEN`
   - `DIGITALOCEAN_ACCESS_TOKEN`

3. Create a GitHub Actions workflow:
```bash
mkdir -p .github/workflows
```

4. The workflow will automatically:
   - Build the Docker image
   - Push to Docker Hub
   - Deploy to Digital Ocean
   - Set up HTTPS with Let's Encrypt

## License
MIT 