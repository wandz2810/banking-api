docker-compose up --build
docker-compose down
curl -X GET http://localhost:8080/customer
curl -X GET http://localhost:8080/customer/123586
docker run -d --name shared_mongodb --network shared-network -p 27017:27017 -v mongo_data:/data/db mongo:7.0.14
