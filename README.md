## Architecture

The following diagram shows the high-level architecture of the system:

<img width="1109" height="551" alt="Lambda drawio" src="https://github.com/user-attachments/assets/02e9687a-e8a6-4d55-8424-41f76b07e470" />

### Architectural Overview

- **Client**  
  Web or mobile clients communicate with the system over HTTPS using JWT-based authentication.

- **Amazon API Gateway (HTTP API + JWT Authorizer)**  
  Acts as the entry point to the system, handling request validation and authorization before invoking backend services.

- **AWS Lambda (NestJS API)**  
  Implements the core business logic in a stateless manner, following REST principles.

- **Amazon DynamoDB (On-Demand)**  
  Stores user data with automatic scaling and minimal operational overhead.

- **Amazon SQS**  
  Decouples synchronous API operations from asynchronous background processing.

- **AWS Lambda Worker**  
  Consumes messages from SQS to process asynchronous tasks independently.

- **Amazon CloudWatch**  
  Collects logs and metrics from API Gateway and Lambda functions for monitoring and observability.

## Key Design Decisions

- **Serverless-first approach** to minimize operational complexity and costs.
- **JWT authorization handled at API Gateway level** to keep Lambda functions focused on business logic.
- **DynamoDB On-Demand capacity** to avoid over-provisioning during low traffic.
- **Event-driven asynchronous processing** using SQS to improve scalability and resilience.
- **Stateless services** to allow horizontal scaling without coordination.
