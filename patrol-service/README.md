# Patrol Service 

## Overview

Patrol service manages the patrolling system. 

There will be two primary user roles for the patrol system: `PATROL` and `DISPATCHER`.

It allows users with `PATROL` role to update near-real-time location.

It allows `DISPATCHER` to assign one existing crime to one user with `PATROL` .

It manages the status of `PATROL` that includes location and availability.

## Functionality

- The service can list all patrol statuses and locations.
- It should be able to get information for a single patrol.
- It should allow `ADMIN` to register an existing user as an officer requiring `userId`, `officerId` and `officerName`.
- It should allow `DISPATCHER` and `ADMIN` to assign patrol to a crime, during which both crime `status` and `patrol_id` should be updated accordingly. 
- It should only allow `PATROL` and `ADMIN` to update their own patrol status, except `ADMIN` has access to all patrol statuses.

## Assumptions

At the user end, patrol is scheduled to update location every second, and every two seconds the location and status will be populated to Postgres for data persistence. 

## Redis

Redis is used for fast updating with patrols. 

### Redis Modelings

- Location updating in Redis does not allow partial updating, only update the full location once to prevent inconsistency