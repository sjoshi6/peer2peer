# Steps to setup project on ubuntu

```
sudo apt-get update
sudo apt-get install golang
```

Create the workspace
```
mkdir go-workspace
cd go-workspace/
mkdir src
mkdir pkg
mkdir bin
```

Set EnV params
```
sudo vi ~/.bashrc

Note: add at the end ; save and restart terminal

export GOPATH=$HOME/go-workspace
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

Clone Project
```
git clone https://github.com/sjoshi6/peer2peer.git
cd go-workspace/src/peer2peer

go get
go install
```

Postgres DB setup
```
sudo apt-get install postgresql postgresql-contrib
sudo -u postgres psql postgres      // set "" as password when asked
sudo -u postgres createdb peer2peer // Creates the DB

```

Start App
```
sudo -u postgres ./peer2peer 8000
```
