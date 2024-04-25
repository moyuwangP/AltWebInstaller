rm -rf alt-web-installer
rm -rf ./dist
cd ./backend
go mod tidy
go build -o ../alt-web-installer main.go
cd ..
cd ./frontend
npm install
npm run build