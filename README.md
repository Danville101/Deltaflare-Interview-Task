# Microservices Architecture for Event Handling and Docker Integration

## Overview
This project implements a microservices architecture using Golang and Docker to handle real-time security events. The system is designed to:

- Generate JSON-formatted security events.
- Store events in InfluxDB.
- Query and display critical events.
- Use NATS for messaging between microservices.

## Architecture
Each microservice is implemented in its own folder, emulating an independent repository structure. This design decision supports the principle of microservices, allowing each service to be developed, deployed, and scaled independently.

### Clean Architecture
To enhance scalability, testability, and modularity, a minimal clean architecture was implemented for each microservice. This architecture separates concerns by dividing the database layer and the service layer, ensuring that each layer can be developed and maintained independently.

### Services

1. **Daemon Service**
   - Continuously generates random security events in JSON format.
   - Publishes events to NATS under the subject `events`.

2. **Client Service**
   - Queries the last 10 events with a criticality level higher than a specified threshold.
   - Displays events in a clear, concise format.
   - Packaged as a Docker container.

3. **Reader Service**
   - Listens for NATS requests to query InfluxDB for specific events.
   - Returns the queried events via NATS.
   - Packaged as a Docker container.

4. **Writer Service**
   - Subscribes to NATS events and writes them to InfluxDB.
   - Packaged as a Docker container.

### Infrastructure
- **NATS**: Facilitates communication between the microservices.
- **InfluxDB**: Stores the security events.

### Docker Compose
Docker Compose is used to manage the microservices network, including the NATS server and InfluxDB container.

## Dependencies

- **github.com/nats-io/nats.go**: NATS Go client for messaging.
- **github.com/joho/godotenv**: Loads environment variables from `.env` files.
- **github.com/influxdata/influxdb-client-go/v2**: InfluxDB client for Go.
- **github.com/stretchr/testify/assert**: Assertion library for unit testing.
- **github.com/stretchr/testify/mock**: Mocking library for unit testing.

## Setup Instructions

1. Clone the repository.
2. Navigate to each service folder and build the Docker images.
3. Use Docker Compose to start the entire system:
   ```sh
   docker-compose up
   ```
4. The services will start, and you can interact with them as needed.

## Running the Services

- The Daemon service will automatically start generating events.
- Use the Client service to query and display events.
- The Reader and Writer services handle data retrieval and storage with InfluxDB.

## Design Choices

- **Independent Services**: Each service in its own folder ensures modularity, making it easier to manage and scale.
- **Clean Architecture**: Implemented to separate concerns between the database and service layers, enhancing scalability, testability, and modularity.
- **Docker Integration**: Containerization simplifies deployment and ensures consistency across environments.
- **NATS for Messaging**: Provides a robust, scalable messaging system for inter-service communication.

## Quality Assurance

- **Testing**: Utilized Testify for unit testing to maintain high code quality and reliability.

## Security

- **Environment Variables**: Critical information is stored in a `.env` file, ensuring sensitive data is not hard-coded.
- **GitHub Ignore**: A `.gitignore` file is used to prevent the `.env` file from being pushed to the public repository, safeguarding sensitive information.
- **Role-Based Access Control (RBAC)**: Implemented least privileged role-based access control for each microservice when connecting to NATS. This ensures that each service has only the necessary permissions required for its operation, thereby minimizing the potential attack surface.
  - **Daemon Service**: Can publish events to the `events` subject.
  - **Writer Service**: Can subscribe to the `events` subject.
  - **Reader Service**: Can publish and subscribe to `query.event`.
  - **Client Service**: Can publish and subscribe to `query.event`.


## Future Updates

Future updates include implementing encryption in transit using TLS and MTLS:

- **TLS**: For securing communication with InfluxDB.
- **MTLS**: For securing communication between NATS and all microservices.

## Conclusion
This microservices setup demonstrates a robust and scalable architecture for handling real-time security events. It leverages Golang for service implementation, Docker for containerization, and NATS and InfluxDB for communication and data storage.

