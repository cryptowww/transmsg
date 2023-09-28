# How to transform message with ruleGo 

## client
it is reponsible for send message to the middlewear server where the message is transformed.

it's not run automatically, so you can active it by send a POST request by postman to th location:
http://localhost:8081/send

## middleWear server

In here, the middleWear(I called it wm) receives message and transform it to the format for the next stop(the server).

## server

It receive the message, and execute other ops.
