# GeoSpatial Extensions for Redis

This tester adds support to the Codecrafter's Redis tester for the GeoSpatial extensions of Redis, allowing users to implement commands like:
- `GEOADD` - Add geospatial data to a key.
- `GEOPOS` - Retrieve the coordinates of members stored using `GEOADD`.
- `GEOHASH` - Calculate the hash of a geospatial key.
- `GEODIST` - Calculate the distance between two members in a geospatial key.

Read about the [stage-wise requirements](./STAGES.md) to see how to implement the commands.

## How to run the tester

To run the tester, you can place your Codecrafter's repository in the folder `/internal/test_helpers/pass_all` and then run `make test_geo_with_redis` to run the appropriate tests.

Please do make sure you do copy your ``spawn_redis_server.sh`` script and the `codecrafters.yml` file as well. 

### Customing the tester

The tester binary runs the tests on the submission and stages marked by the follwowing environment variables:
- `CODECRAFTERS_SUBMISSION_DIR` - the root directory of your code submission.
- `CODECRAFTERS_TEST_CASES_JSON` - the test cases in JSON format.

You can set them as you want to change the defaults. It would be most convenient to copy the appropiate testing script from the [Makefile](./Makefile) and then change the environment variables to your needs.

