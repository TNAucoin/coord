## Project Layout

- Application folder: This folder contains microservice business logic, 
which is a combination of the domain model that refers to a business entity and an API that exposes core functionalities to other modules.

- Port folder: This folder contains contract information for integration between the core application and third parties.
This can be a contract about accessing core application features or about specifying available features for a database system, if one is used for the persistence layer.

- Adapter folder: This folder contains concrete implementation for using contracts that are defined in ports. 
For example, gRPC can be an adapter with a concrete implementation that handles requests and uses an API port to access core functionalities, 
such as if you have an application with some functionalities and will expose it to customers. 

