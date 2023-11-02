# chatapp

Ive started over on this project in another repository.
Reasons:
1. I've decided to rethink the Friends system of the chat application. A DB table for friend requests seems necessary.
2. I want to swap out mongoDB for a relational Database where i can write the SQL queries myself instead of relying on the Bson document format with aggregation pipelines etc. 
3. I want to approach the websocket handling differently. I'll go for a process where several distributionhubs are launched to enable users to have their separate conversations active in those. 
