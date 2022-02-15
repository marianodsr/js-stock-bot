# js-stock-bot

 Stock bot for jobsity coding challenge. Meant to be run alongside https://github.com/marianodsr/js-chat-room

 ## How to run  
 
 Run docker-compose up by being at the root of the repository
  
 This application is decoupled from chat room, but docker compose in chat-room is the one in charge of spinning  
 the rabbitMQ container  
 
 If for some reason you want to start this application by itself, then you are going to need to go to chat-room and run:
 
 docker-compose run rabbitMQ-container

since this is a dependency.
 
 Clone this repo and follow https://github.com/marianodsr/js-chat-room readme.
 
 Happy testing!
