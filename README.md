# location_store

clientID = mWQQAztWIqpmoqbFWmcbQA
apikey = Z5QzsWaC1oiBNg-vf9ZylfBdphfGkbqYr7B1FCmvtn5eDwJf5Ql1ANidgc5LdgruYWKjoI56LG9eg4za7muZ5kh3gyuPSFMEr9q2h0kDdQQo3q4N9FcKqCVI3IwjZHYx

Build an API that provides a way for the user to upload locations by using a .csv file and then populate the locations with metadata from the yelp API.

Using

Golang

Postgres + Postgis

use a worker pool with a max of 5 concurrent workers to manage the number of concurrent goroutines

Test data from places_interview.csv

Requirements

The queue system has to be built from scratch

If you use external libraries, you have to justify why they are needed. No library for the queue is allowed.

Acceptance Criteria

An endpoint that accepts a .csv file with columns: lat, long, place_id

Store the new entries in the database

Don’t save duplicated places

Check which location don’t have near by suggestions

In a background process, retrieve near by locations in the background and store them in the database

Don’t save duplicated suggestions

An endpoint to list all the locations

Make a list of all the tests cases that you would need to run but implement only the 3 most important ones



YELP documentation
https://api.yelp.com/v3/businesses/search?term=delis&latitude=37.786882&longitude=-122.399972

https://docs.developer.yelp.com/docs/fusion-intro

![tree](https://user-images.githubusercontent.com/46742766/232639719-83a0adb1-d966-449e-853d-277e6d137cbc.png)
![api1](https://user-images.githubusercontent.com/46742766/232639732-11904016-bcfa-4e17-a4f1-500bedb11051.png)
![api2](https://user-images.githubusercontent.com/46742766/232639738-bb8ac424-809e-44f2-bd99-85a4e43a151f.png)
![test](https://user-images.githubusercontent.com/46742766/232639742-95455a33-08aa-4d93-8151-c335d1160498.png)

