# QULER
Read message from RabbitMQ cross platform  
1. Download app from [releases](https://github.com/vmpartner/quler/releases)  
2. Change config ```app.conf```
3. Run ``` ./quler_windows_64.exe ``` 

# Params
#### App
**name** - Name of app, it's shown in rabbit connections  
**sync_each_message** - Each message will be written to disk immediately  

#### MQ
**queue_source** = Read queue
**ack_message** - Ack message after received
**limit_messages** - Stop read queue after N messages

#### File
**message_per_file** - Each message in separate file   
**path** - Path and mask for files where % is number of message = ./result/mess_%.txt
