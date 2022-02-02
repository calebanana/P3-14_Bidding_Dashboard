# ETI Assignment 2
## Done by 
Caleb Goh En Yu

---
## 1. Introduction
This project was created as part of an assignment for the Emerging Trends in IT (ETI) module in my diploma course. It involves working with fellow peers to design, implement and containerise several microservices and REST APIs to bring about various functions of a simulated online institution learning management system, named EduFi. This online platform was designed with many different features, and the set of features that I have taken up would be the class bidding system, as described below.

**Package 3.14: Bidding Dashboard**
- 3.14.1.	List all classes
- 3.14.2.	Search for classes
- 3.14.3.	View class information
- 3.14.4.	Bid classes with ETI tokens
- 3.14.5.	View, update, delete own bids
- 3.14.6.	View anonymized bids of classes
- 3.14.7.	View all bids of classes
- 3.14.8.	Show results of bidding

---
## 2. Design Consideration of Microservices
From all the requirements stated above, most of them would require a user interface where information can be displayed to users (students) while also allowing them to enter their inputs. As such, a dedicated microservice for a front-end web application is very much needed.

To handle all the logic and perform the stated functions, another microservice existing in the form of an API must be present. This API should be able to accept different types of requests depending on the type of action that the users are performing. From the requirements above, each one of them would fit into one of the four types of API request types: ```GET```, ```PUT```, ```POST``` and ```DELETE```. Detailed descriptions of the available API endpoints can be found [here](./microservices.md).

Lastly, a database has to be set up in order to provide a persistent platform where bidding data can be stored. As the overall structure of the information that I am working with is very nested, a document storage database would be the most suitable. As such, I have decided to implement the bidding database using MongoDB, with the following structure.

```
{
  "SemesterStartDate":"...",
  "SemesterModules":[
    {
      "ModuleCode":"...",
      "ModuleName":"...",
      "ModuleClasses":[
        {
          "ClassCode":"...",
          "ClassBids":[
            {
              "StudentID":"...",
              "BidAmt":0,
              "BidStatus":"..."
            },
            ...
          ]
        },
        ...
      ]
    },
    ...
  ]
}
```

---
## 3. Architecture Diagram
The structure of the microservices implemented is as shown in the diagram below. It consists of a web application microservice which provides the user interface and serves as the front end. The web application communicates with the Bidding API microservice implemented using GO through HTTP requests made to the API endpoint. . 

Apart from handling all the logic, the Bidding API is able to create, update and remove records from the bidding database implemented with MongoDB and run on a separate microservice.

![Architecture Diagram](/architecture_diagram.png)
---
## 4. Set-Up Instructions
### Deployment on Class Server
1. On the server, pull the three Docker images for the project from Docker Hub.

    ```
    docker pull calebanana/bidding_db
    ```
    ```
    docker pull calebanana/bidding_dashboard
    ```
    ```
    docker pull calebanana/bidding_api
    ```
2. Run the containers.

    ```
    docker run -d calebanana/bidding_db
    ```
    ```
    docker run -d calebanana/bidding_dashboard
    ```
    ```
    docker run -d calebanana/bidding_api
    ```
3. The Front-End page will then be accessible on http://localhost:8220/biddingDashboard