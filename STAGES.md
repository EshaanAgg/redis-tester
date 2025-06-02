# Stage 1: Respond to GEOADD command
  
- **Slug**: `kj2`
- **Difficulty**: `easy` 
- **Marketing**: In this stage, you will implement the `GEOADD` command to add geospatial data to Redis.

In this stage, you will implement the support for the [GEOADD](https://redis.io/docs/latest/commands/geoadd/) command. 

Geospatial data for multiple locations is stored under a single key, and each location is represented by a longitude, latitude, and name. All the locations under a key are also known as the "members" of that key.

The `GEOADD` command stores the data internally, and returns a single integer value indicating the number of members added to the key.

In this stage, the tester would always send only 1 member to the `GEOADD` command, so you can cut corners and hard code your response to always return `1` when the command is `GEOADD`.

### Tests

The tester will execute your program like this:

```bash
$ ./spawn_redis_server.sh
```

It will then add a single member to the key `cities` with the following command:

```bash
$ redis-cli GEOADD cities  -0.127758 51.507351 london
```

The expected output is:

```
:1\r\n
```

which is the RESP2 encoding for the [simple integer](https://redis.io/docs/latest/commands/protocol/#simple-integer) `1`.


# Stage 2: Handle invalid GEOADD coordinates

- **Slug**: `pq5`
- **Difficulty**: `easy`
- **Marketing**: In this stage, you will improve your `GEOADD` implementation by validating the latitude and longitude values provided by the user.

In this stage, you will add validation logic to your `GEOADD` implementation. Redis requires longitude and latitude values to fall within a valid range:

- Longitude must be between -180 and 180 (inclusive).
- Latitude must be between -85.05112878 and 85.05112878 (inclusive).

If a location is provided with invalid coordinates, your Redis server must return an error.

The expected error message is:

```
-ERR invalid longitude,latitude pair <longitude>,<latitude>\r\n
```

where `<longitude>` and `<latitude>` should be printed with six decimal places.

### Tests

The tester will execute your program like this:

```bash
$ ./spawn_redis_server.sh
```

It will then send various `GEOADD` commands to test your validation, such as:

```bash
$ redis-cli GEOADD cities -180.000123 12.345678 invalid_long
```

Expected output:

```
-ERR invalid longitude,latitude pair -180.000123,12.345678\r\n
```

If the coordinates are valid, your Redis implementation should continue to return `:1\r\n` just like in the previous stage.


# Stage 3: Support multiple entries in GEOADD

- **Slug**: `si2`
- **Difficulty**: `easy`
- **Marketing**: Redis's `GEOADD` command supports adding multiple geospatial items in one go. In this stage, you'll implement this batch insert capability with proper validation.

Up until now, your `GEOADD` implementation may have only supported inserting one location per command. Redis allows clients to insert multiple locations in a single `GEOADD` command, and you will now need to support this feature.

Here’s an example of such a command:

```
GEOADD cities 13.361389 38.115556 Palermo 15.087269 37.502669 Catania
```

This adds two members with their respective coordinates to the `cities` key.

Just like in the previous stage, **all** coordinates must be valid. If **any** one of the provided coordinates is invalid:

- Do not add **any** of the elements.
- Return the error message:

```
-ERR invalid longitude,latitude pair <longitude>,<latitude>\r\n
```

### Tests

The tester will execute commands like:

```bash
$ redis-cli GEOADD cities_1 13.361389 38.115556 Palermo 15.087269 37.502669 Catania
````

Expected response:

```
:2\r\n
```

It will also test mixed valid and invalid inputs, like:

```bash
$ redis-cli GEOADD cities_2 12.34 56.78 ValidCity 91.123456 181.123456 InvalidCity
```

Expected response:

```
-ERR invalid longitude,latitude pair 91.123456,181.123456\r\n
```


# Stage 4: Implement GEOPOS command

- **Slug**: `rt7`
- **Difficulty**: `easy`
- **Marketing**: In this stage, you will implement the `GEOPOS` command to retrieve the coordinates of members stored using `GEOADD`.

In this stage, you will implement support for the [GEOPOS](https://redis.io/docs/latest/commands/geopos/) command.

The `GEOPOS` command returns the longitude and latitude of one or more members of a geospatial key. The coordinates are returned as RESP arrays of bulk strings.

For example, if you have already added the member `london` under the key `cities`, and then you query for the same using the following command:

```bash
$ redis-cli GEOPOS cities london
```

you would get:

```
*1\r\n*2\r\n$10\r\n-0.127758\r\n$9\r\n51.507351\r\n
```

which is equivalent to RESP encoding for:

```
[
  [
    "-0.127758",
    "51.507351"
  ]
]
```

In this stage, you will only need to retrieve a single member which is guranteed to exist.

### Tests

The tester will run your program like this:

```bash
$ ./spawn_redis_server.sh
```

It will then execute the `GEOADD` command to add a member, and then run the `GEOPOS` command to retrieve its coordinates:

```bash
$ redis-cli GEOADD cities -0.127758 51.507351 london
$ redis-cli GEOPOS cities london
```

The expected output is:

```
*1\r\n*2\r\n$10\r\n-0.127758\r\n$9\r\n51.507351\r\n
```

# Stage 5: Handle non-existing members in GEOPOS

- **Slug**: `uv3`
- **Difficulty**: `easy`
- **Marketing**: In this stage, you will improve your `GEOPOS` implementation to handle cases where the requested member does not exist.

In this stage, you will enhance your `GEOPOS` implementation to handle cases where the requested member does not exist in the geospatial key.

If a member does not exist, the `GEOPOS` command should return a RESP array with a single `nil` value. For example, if you run:

```bash
$ redis-cli GEOPOS cities non_existing_member
```

The expected output should be:

```
*1\r\n$-1\r\n
```

### Tests

The tester will execute your program like this:

```bash
$ ./spawn_redis_server.sh
```

It will then run the `GEOPOS` command for some non-existing member:

```bash
$ redis-cli GEOPOS cities non_existing_member
```

The expected output is:

```
*1\r\n$-1\r\n
```

It will also test existing members to ensure they return the correct coordinates like the previous stage.

# Stage 6: Support multiple members in GEOPOS

- **Slug**: `wx2`
- **Difficulty**: `easy`
- **Marketing**: In this stage, you will extend your `GEOPOS` implementation to support retrieving coordinates for multiple members in a single command.

In this stage, you will enhance your `GEOPOS` implementation to support retrieving coordinates for multiple members in a single command.

The `GEOPOS` command can accept multiple member names, and it should return their coordinates in a single response. The response should be a RESP array containing arrays of coordinates for each member.

For example, if you have already added the members `london` and `paris` to the key `cities`, and then you query both of their position together using the command:

```bash
$ redis-cli GEOPOS cities london paris
```

you would get back:

```
*2\r\n*2\r\n$10\r\n-0.127758\r\n$9\r\n51.507351\r\n*2\r\n$7\r\n2.352222\r\n$9\r\n48.856613\r\n
```

which is equivalent to RESP encoding for:

```
[
  [
    "-0.127758",
    "51.507351"
  ],
  [
    "2.352222",
    "48.856613"
  ]
]
```

Each element of the array corresponds to one of the queries location. The first element of each array element in the longitude, and the second is the lattitude.

### Tests

The tester will execute your program like this:

```bash
$ ./spawn_redis_server.sh
```

It will then run the `GEOADD` command to add multiple members, and then execute the `GEOPOS` command to retrieve their coordinates:

```bash
$ redis-cli GEOADD cities -0.127758 51.507351 london 2.352222 48.856613 paris
$ redis-cli GEOPOS cities london paris
```

It is expected to return the correct coordinates for both members in the format described above. Your implementation should handle both existing and non-existing members gracefully, returning `nil` for any non-existing members.


# Stage 7: Add GEOHASH command

- **Slug**: `xy3`
- **Difficulty**: `hard`
- **Marketing**: In this stage, you’ll implement the `GEOHASH` command, which returns geohash strings for specified members in a geospatial index.

In this stage, you'll implement support for the [`GEOHASH`](https://redis.io/docs/latest/commands/geohash/) command.

The `GEOHASH` command returns geohash-encoded strings representing the positions of specified members of a geospatial key. These strings are 11-character base32 encodings of the geospatial coordinates. Redis represents positions of elements using a variation of the Geohash technique where positions are encoded using 52 bit integers. You can read more about the encoding and decoding of geohashes on [Wikipedia](https://en.wikipedia.org/wiki/Geohash).

For example, you can get the geohash of the member `london` under the key `city` by running:

```bash
$ redis-cli GEOHASH cities london
```

The output of the same would be:

```
*1\r\n$11\r\gcpvj0dup10\r\n
```

which is the RESP2 representation of an array with single bulk string `gcpvj0dup10`.

The tester would not expect you to return the exact same geohash as Redis does, but rather a geohash which when decoded gives the longitude and latitude of the member with a precision of 5 decimal places. 

### Tests

The tester will run your program like this:

```bash
$ ./spawn_redis_server.sh
```

It will then execute the `GEOADD` command to add a member, and then run the `GEOHASH` command to retrieve its geohash:

```bash
$ redis-cli GEOADD cities -0.127758 51.507351 london
$ redis-cli GEOHASH cities london
```

The expected output is:

```
*1\r\n$11\r\gcpvj0dup10\r\n
```

or similar, as long as the geohash decodes to the correct longitude and latitude with a precision of 5 decimal places.

# Stage 8: Support mutliple memebers in GEOHASH command

- **Slug**: `zw4`
- **Difficulty**: `easy`
- **Marketing**: In this stage, you will extend your `GEOHASH` implementation to support multiple members in a single command.

In this stage, you will enhance your `GEOHASH` implementation to support retrieving geohashes for multiple members in a single command.

The `GEOHASH` command can accept multiple member names, and it should return their geohashes in a single response. The response should be a RESP array containing bulk strings for each member's geohash.

For example, if you want to query for the geohash of the members `london` and `paris` under the `cities` key, you can run the command:

```bash
$ redis-cli GEOHASH cities london paris
```

It would return:

```
*2\r\n$11\r\gcpvj0dup10\r\n$11\r\u09tvw0f6s0\r\n
```

which is the RESP2 encoding for:

```
[
  "gcpvj0dup10",
  "u09tvw0f6s0"
]
```

### Tests

The tester will execute your program like this:

```bash
$ ./spawn_redis_server.sh
```

It will then run the `GEOADD` command to add multiple members, and then execute the `GEOHASH` command to retrieve their geohashes:

```bash
$ redis-cli GEOADD cities -0.127758 51.507351 london 2.352222 48.856613 paris
$ redis-cli GEOHASH cities london paris
```

The expected output is:

```
*2\r\n$11\r\gcpvj0dup10\r\n$11\r\u09tvw0f6s0\r\n
```

or similar, as long as the geohashes decode to the correct longitude and latitude with a precision of 5 decimal places for each member.

# Stage 9: Add GEODIST command

- **Slug**: `qv5`
- **Difficulty**: `medium`
- **Marketing**: In this stage, you will implement the `GEODIST` command to calculate the distance between two members in a geospatial index.

In this stage, you will implement support for the [`GEODIST`](https://redis.io/docs/latest/commands/geodist/) command.

The `GEODIST` command calculates the distance between two members of a geospatial key. The command supports multiple units for the distance measurement, but defaults to meters if no unit is specified. 

In this stage, you would only need to handle the default case and return the distance betweent the two provided members in meters. Redis assumes that the Earth is a perfect sphere while calculating the distance, so you can use [the Haversine formula](https://en.wikipedia.org/wiki/Haversine_formula) for this calculation.

You can run the followng command to get the distance between the members `london` and `paris`, stored under the key `cities`:

```bash
$ redis-cli GEODIST cities london paris
```

It would return a bulk string with the distance in meters (`343646.2197`), like:

```
$11\r\n343646.2197\r\n
```

### Tests

The tester will execute your program like this:

```bash
$ ./spawn_redis_server.sh
```

It will then run the `GEOADD` command to add two members, and then execute the `GEODIST` command to calculate the distance between them:

```bash
$ redis-cli GEOADD cities -0.127758 51.507351 london 2.352222 48.856613 paris
$ redis-cli GEODIST cities london paris
```

The expected output is:

```
$11\r\n343646.2197\r\n
```

The tester would not expect you to return the exact same distance as Redis does, but would accept any distance which is within a 1% margin of error from the correct distance calculated by the Haversine formula.

# Stage 10: Support multiple units in GEODIST 

- **Slug**: `qv8`
- **Difficulty**: `easy`
- **Marketing**: In this stage, you will extend your `GEODIST` implementation to support multiple units for distance measurement.

In this stage, you will enhance your `GEODIST` implementation to support multiple units for distance measurement.

The `GEODIST` command can accept an optional unit argument to specify the desired unit for the distance measurement. The supported units are:
- `m` for meters (default)
- `km` for kilometers
- `ft` for feet
- `mi` for miles

If a unit is provided, your implementation should convert the distance from meters to the specified unit before returning it. For example, to get the distance between the members `london` and `paris` in the key `cities` in kilometers, you can run the command:

```bash
$ redis-cli GEODIST cities london paris km
```

It would return a bulk string with the distance in kilometers (`343.6462197`), like:

```
$11\r\n343.6462197\r\n
```

## Tests

The tester will execute your program like this:

```bash
$ ./spawn_redis_server.sh
```

It will then run the `GEOADD` command to add two members, and then execute the `GEODIST` command with a unit to calculate the distance between them:

```bash
$ redis-cli GEOADD cities -0.127758 51.507351 london 2.352222 48.856613 paris
$ redis-cli GEODIST cities london paris km
```

The expected output is:

```
$11\r\n343.6462197\r\n
```

The tester would not expect you to return the exact same distance as Redis does, but would accept any distance which is within a 1% margin of error from the correct distance in the appropiate unit.
