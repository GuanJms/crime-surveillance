# Crime Service

Crime service should enable gRPC listening at 50001.

## Features
- Crime service should be able to receive request of listing all the crimes.
- it can connect to postgres and execute query.


## Known issues

- dbTimeout issue : Setting fixed dbTimeout with 3 seconds would be troublesome when it comes to listing all crimes when base is large
    -  solutions: 
        - break up the listing all crimes to small batches
        - set up a dynamic dbTimeout based on the possible estimation