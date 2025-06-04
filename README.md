# GeoSpatial Extension for Redis

This repository modifies the [Codecrafters Redis tester](https://github.com/codecrafters-io/redis-tester) to include support for the GeoSpatial extension of Redis. The GeoSpatial extension allows you to store, query, and manipulate geospatial data, such as locations and distances, directly within Redis. 

This extension currently adds 10 new stages, in which you will implement the following commands:

- `GEOADD` - Add geospatial data to a key.
- `GEOPOS` - Retrieve the coordinates of members stored using `GEOADD`.
- `GEOHASH` - Calculate the hash of a geospatial key.
- `GEODIST` - Calculate the distance between two members in a geospatial key.

Read about the [stage-wise requirements](./STAGES.md) to see how to implement the commands.

## How to run the tester

To run the tester, you can place your Codecrafter's repository in the folder `/internal/test_helpers/pass_all` and then run `make test_geo_with_redis` to run the all the Geospatial stages. 

Please do make sure you do copy your `spawn_redis_server.sh` script and the `codecrafters.yml` file as well (and delete the current ones)!

### Customing the tester

The tester binary runs the tests on the submission and stages marked by the follwowing environment variables:
- `CODECRAFTERS_SUBMISSION_DIR` - the root directory of your code submission.
- `CODECRAFTERS_TEST_CASES_JSON` - the test cases in JSON format.

If you do not want to copy your Codecrafter's repository to the `/internal/test_helpers/pass_all` folder, you can instead modify the [`Makefile`](./Makefile)'s `test_geo_with_redis` command, and change the `CODECRAFTERS_SUBMISSION_DIR` to point to your Codecrafter's repository! You can change the `CODECRAFTERS_TEST_CASES_JSON` to run specific stages, or use the other make commands to run tests for other stages.


