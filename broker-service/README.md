# Broker Service
Broker service is the gateway to all microservices - HTTP request will come in and pass in and trigger response on corresponding services.

## Crime broker service

### Features

- Crime broker service should be able to request crimes and receive the crime responses from crime server. 
    - During requesting, crimb broker can request with filters that includes crime_id, user_id, street, city, state, date, status, and keywords for description filtering.
- Crime broker service should be able to submit a new crime with reporter_id (user_id), description, street, city, state, latitude, and longitude.
- Crimbe broker service should be able to update a new crime 