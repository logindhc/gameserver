@echo off

echo build go protocol file...
cd internal
cd protocol
protoc --go_out=./../pb/ --go_opt=paths=source_relative *.proto
cd ../..
echo build go proto complete!