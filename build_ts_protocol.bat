@echo off

if exist outts (
  rd /s /q outts
)
md outts

protoc --proto_path=internal/protocol/ --ts_proto_out=outts/ internal/protocol/*.proto

echo build ts proto complete!
