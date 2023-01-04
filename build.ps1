$dir = Get-Item .\
cd src
go mod tidy
go build -o ..\build\diceroller.exe .\
cd $dir
Copy-Item -Force -Recurse .\js .\build\
Copy-Item -Force -Recurse .\css .\build\
Copy-Item -Force -Recurse .\res .\build\
Copy-Item -Force -Recurse .\db .\build\
Copy-Item -Force -Recurse .\templates .\build\
Copy-Item -Force -Recurse .\diceroller.conf .\build\