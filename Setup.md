Setup Instructions
-----
1. Run cmd `docker build -t mycli:latest .` to build the docker image for this application
2. Run cmd `docker-compose up` to run an instance of MySQL as a container

These commands will have the basic setup ready.

Populate DB
----
This application reads the given input json file and stores its contents in the MySQL instance running
after running the above commands. This is intentional. By doing this, I am making read path faster over
the write path.

Json fixtures file with recipe data. Download [Link](https://raw.githubusercontent.com/hellofreshdevtests/json-fixtures/content/hf_test_calculation_fixtures.tar.gz).

This input file is more than 1GB in size and has more than 7Mn records. Inserting this data into db
takes good 30 minutes. Hence, stop the execution of application if this time period is too long to wait.
You can work with partial data and see the expected output too. Run below command to insert data.
1. ` docker run --rm -v <host_directory>:/data mycli:latest --input=/data/<file_name>.json`
2. Run above command with different files names to insert more data.

Replace placeholders `host_directory` and `file_name` with values applicable to you

Flags
---
This application supports custom flags. These can be used during execution to slice and dice data as per your needs.
If custom values are not provided during execution, then application assumes some reasonable defaults. These have been
listed in README.md file.

1. Run `docker run  mycli:latest  --help`

- `-input string`
    - **Description**: Path to the JSON file.
    - **Usage**: Specify the path to the JSON file that contains the input data.

- `-postcode string`
    - **Description**: Custom Postcode.
    - **Usage**: Provide a specific postcode to filter or use in the report.

- `-fromTime string`
    - **Description**: Custom From Time.
    - **Usage**: Specify the starting time for the report or filter. Follow strictly 12 HOUR time format. e.g 3PM, 9PM
- `-toTime string`
    - **Description**: Custom To Time.
    - **Usage**: Specify the ending time for the report or filter. Follow strictly 12 HOUR time format. e.g 1AM, 12AM

- `-recipe string`
    - **Description**: Recipe Names to search for. Provide this as a comma-separated list with no spaces in between.
    - **Usage**: List the recipe names you want to search for in the input data.

Custom Command Examples
-----
- `docker run  mycli:latest .`
    - **Description**: This will return data with defaults values selected for all the flags
- `docker run  mycli:latest -postcode=10120`
    - **Description**: This will return data with for the selected value of postcode
- `docker run  mycli:latest --fromTime=7AM --toTime=5PM`
    - **Description**: This will return data between given from time and to time
- `docker run  mycli:latest  --fromTime=7AM --toTime=5PM --postcode=10121`
    - **Description**: This will return data between given from time, to time and postcode

Output Format
----
All commands return data in format as mentioned in README.md file

Note
-----------
This is `NOT` a production grade application. This stills needs some improvements