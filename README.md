# Place Store Example

## Project Tree
![tree](https://user-images.githubusercontent.com/46742766/232639719-83a0adb1-d966-449e-853d-277e6d137cbc.png)

## APIs
I didn't dockerize or deployed the project anywhere.


But If you ask I can prepare dockerfile or 


provide you an IP for accessing to services.

### First API:
![api1](https://user-images.githubusercontent.com/46742766/232639732-11904016-bcfa-4e17-a4f1-500bedb11051.png)

### Second API:
![api2](https://user-images.githubusercontent.com/46742766/232639738-bb8ac424-809e-44f2-bd99-85a4e43a151f.png)

## Tests
![test](https://user-images.githubusercontent.com/46742766/232639742-95455a33-08aa-4d93-8151-c335d1160498.png)

### Additional test cases should be tested:
1. What if the number of places in the csv file is less than the maxWorker?

    Does it cause an error?

2. How to act if there is an error like too mant connections to DB?

    I implemented pool connection, but does it enough?

3. What if there is an error on one of the jobs?

	Should we stop the whole process?

    Should we continue the process?
	
    Should we retry the failed job?
	
    Should we retry the whole process?

 4. Yelp says that Limit + Offset must be less than or equal to 1000.

    This avoids me to get all nearby places.

5. In case of requests are too fast like exceeding rate limiting to the Yelp API, or any other error, we should implement some retry mechanism until it get success.

    For "rate limiting error" recursion is applied, but not for any kind of errors.

6. The amount of unique nearby places coming from Yelp changes. What does caues it?

    As I recognized the nearby places coming from Yelp cannot be ordered. Does the order of places change? 

I also placed these as comments in the code.

## Further Improvements

Implementing concurrent jobs for 'searching nearby places on Yelp and inserting them to DB', decreased reponse time from 1.5-2 mins to 30-40 secs.

Searching on Yelp concurently with different offset parameters for each individual csv place can cause improvements on reponse time.

## Further Questions
I enabled the postgis extension as you can see in functions. 

As I asked before "I may ask questions about the fields on postgis later, the meaning of fields etc." I was trying to ask you meaning of these fields.
![1](https://user-images.githubusercontent.com/46742766/232644701-590af9f3-5687-4f8b-82b1-ae226bcb7ed0.png)

I enabled but not used anything about postgis.
![2](https://user-images.githubusercontent.com/46742766/232644722-b749415d-bd6a-4c45-8bf8-8f7583aa084b.png)
