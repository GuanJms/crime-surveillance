# Patrol Service 

## Overview

The Patrol service manages the patrolling system. 

There are two primary user roles for the patrol system: `PATROL` and `DISPATCHER`.

It allows users with the `PATROL` role to update their near-real-time location.

It allows `DISPATCHER` to assign an existing crime to a user with the `PATROL` role.

It manages the status of `PATROL` users, which includes location and availability.

## Functionality

- The service can list all patrol statuses and locations.
- It should be able to get information for a single patrol.
- It should allow `ADMIN` to register an existing user as an officer, requiring `userId`, `officerId`, and `officerName`.
- It should allow `DISPATCHER` and `ADMIN` to assign patrol to a crime, during which both the crime `status` and `patrol_id` should be updated accordingly. 
- It should only allow `PATROL` and `ADMIN` to update their own patrol status, except `ADMIN` has access to all patrol statuses.

## Assumptions

At the user end, patrol is scheduled to update location every second, and every two seconds the location and status will be populated to Postgres for data persistence. 

## Redis

Redis is used for fast updating with patrols. 

### Redis Modeling

- Location updating in Redis does not allow partial updates; only the full location can be updated once to prevent inconsistency.

### Persist Database

The system records the updated patrolID into a Redis set and pushes the information into the Persist Database (Postgres database). A Lua script is executed for its atomicity during the syncing process for reading and immediately removing the set of userIDs in the dirty key set (`dirty_keys_patrol`). 