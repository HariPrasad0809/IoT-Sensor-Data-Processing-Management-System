# IoT-Sensor-Data-Processing-Management-System

IoT Sensor Data Processing & Management System using Go (Echo), MySQL & microservices.  
- Microservice A generates sensor data with configurable frequency.  
- Microservice B receives, stores & serves data via REST APIs (with auth, filtering, pagination).  
- Includes Docker, Clean Architecture, ERD, diagrams & API docs.  

## Docker Hub Images

- [Microservice A on Docker Hub](https://hub.docker.com/r/hari0123456789/microservice-a)  
- [Microservice B on Docker Hub](https://hub.docker.com/r/hari0123456789/microservice-b)  

### Pull Images
```bash
docker pull hari0123456789/microservice-a:latest
docker pull hari0123456789/microservice-b:latest
