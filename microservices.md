# Bidding API Reference
The Bidding API is hosted on Port 8221, and it accepts four types of requests as described in the table below. For each of the requests, it takes in the semester start date as a parameter.

For the ```GET``` function, three optional parameters can be passed through a query string: ```classCode```, ```studentId``` and ```filtered```.

For the ```POST``` function, it accepts either zero or two query string parameters: ```classCode```, ```studentId```.

For the ```PUT``` and ```DELETE``` functions, the ```classCode```, ```studentId``` parameters have to be passed in through a query string.


|                                                API Request                                               |                            Description                           |
|:--------------------------------------------------------------------------------------------------------:|:----------------------------------------------------------------:|
| GET /api/v1/bid/:semStartDate                                                                            | Get all bids of all classes                                      |
| GET /api/v1/bid/:semStartDate<br>               ?classCode=...                                           | Get all bids of specific class                                   |
| GET /api/v1/bid/:semStartDate<br>               ?studentId=...                                           | Get all classes (Only returns<br>bids made by specified student) |
| GET /api/v1/bid/:semStartDate<br>               ?studentId=...<br>               &filtered=true          | Get all classes that a<br>specific student has bid for           |
| POST /api/v1/bid/:semStartDate                                                                           | Create empty classes and modules                                 |
| POST /api/v1/bid/:semStartDate<br>                ?classCode=...<br>                &studentId=...       | Add student bid for<br>specific class                            |
| PUT /api/v1/bid/:semStartDate<br>               ?classCode=...<br>               &studentId=...          | Edit student bid amount<br>for specific class                    |
| DELETE /api/v1/bid/:semStartDate<br>                  ?classCode=...<br>                  &studentId=... | Delete student bid<br>for specific class                         |

---
[Back to main README](./README.md)