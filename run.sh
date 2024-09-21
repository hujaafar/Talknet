docker build -t forum .
docker container run -p 8080:8080  --name forumContainer forum
docker stop conp