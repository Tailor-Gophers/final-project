##  Alidada Features

- Using a layered architecture for higher readability and easier changes.
- Ali Dada generally includes the features of registration, login, logout, adding a new passenger, booking tickets, paying and viewing tickets.
- Using the cache layer to store the information of the planes that arrived from the MAC system.
- Ali Dada has the feature of floating price or smart price, which decides at what price to sell a ticket based on the remaining capacity of an airplane.
- Ali Dada includes a smart sorting feature that uses the capacity of the planes to decide which plane to display higher.
- Unit tests and integration tests have been written for this system.
- It has different flight classes and different prices, as well as different cancellation conditions for each of them, depending on the time, there may be different penalties for canceling the ticket.
- With the possibility of issuing tickets at the airport and assigning flight seats with limited access to the admin.

# How to run

## Run with Docker

To run the project, you only need to install Docker on your system


    docker compose up -d

An explanation of the function of each container:



The name of the container  | Function
------------- | -------------
AliDada   | Main system
MockApi  | Simulation of the central system
Qsms  | Main system
mysql1  | AliDada DB
mysql2  | MockApi DB
mysql3  | Qsms DB
redis  | Caching system

## Model diagram
![diagram](alidada/static/diagram.png )

We should note that Flight and FlightClass are in the MockAPI system and are shown in this section only to show the dependencies.
## Caching scenario
It is like a remember scenario, and if there is no key in the cache, it receives it from the main source and stores it to respond from the cache in case of repetition. Of course, it does not read vital information such as the capacity from the cache.


                    
![diagram](alidada/static/senario.png )
## AliDada Postman Document
You can see the document of all APIs along with their examples in the collection below for Ali Dada
[AliDada Postman](https://documenter.getpostman.com/view/16800432/2s93zCYLT1 "AliDada Postman")


# Qsms by Tailor Gophers

### Features

- User and admin registeration and authorization.
- User wallet and payments.
- Creating and Contacts.
- Creating and editing Phonebooks.
- Buying a number from available numbers.
- Renting a number for a monthly basis.
- Sending messages to numbers, contacts and phonebooks.
- Bad words filter.
- Creating and using message templates.
- Setting message scheduler for sending periodic messages with spedified interval.
-  Full admin controll such as suspending and un suspending user, searching messages, counting user messages.


###Description
- This project is developed over MVC(Model View Controller) patter where requests are handled and validated by the Controller layer and corresponding functions will be called from this layer.
- This project is using MySQL for persisting data and Gorm for Object Reletional Mapping.
- you can see the database schema in the diagram below
![](qsms/static/qsms_db_scheme.jpg)

###How to use
- For using the app you must send http requsts after running it you can see full guide to how to send requests in this [postman collection](https://documenter.getpostman.com/view/15181898/2s946fcsBn).
