::protoc --go_out=. account_data.proto
::protoc --go_out=. globalItem_data.proto
::protoc --go_out=. guild_data.proto
::protoc --go_out=. player_data.proto
::protoc --go_out=. ranking_data.proto
::protoc --go_out=. role_data.proto
::protoc --go_out=. net_struct.proto


::protoc -I=$SRC_DIR --python_out=$DST_DIR addressbook.proto

protoc --go_out=. node.proto
::protoc --python_out=. msg.proto
protoc --go_out=. store.proto
protoc --go_out=. message.proto

pause