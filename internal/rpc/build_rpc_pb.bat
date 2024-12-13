@echo off

echo build rpc file...
protoc --go_out=./ --go_opt=paths=source_relative *.proto
echo build rpc proto complete!