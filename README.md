# Prepaid Card

This is a development exercise for building a prepaid card service written in Go.

## Prerequisites

* [Docker](https://www.docker.com/)
* [Git](https://git-scm.com/)
* [GNU Make](https://www.gnu.org/software/make/) *(optional)* - Make is used for convenience as a shortcut to commands.
If you prefer, you can always run the commands, which are listed in the [Makefile](Makefile).


## Installation

Follow the steps to setup your development workspace:
1. Clone the repository
    ```bash
    $ git clone git@github.com:sepetrov/prepaidcard.git ~/Projects/prepaidcard
    ```
1. cd into the project directory
    ```bash
    $ cd ~/Projects/prepaidcard
    ```
1. Create file with environment variables for Docker from the template `.env.dist`
    ```bash
    $ cp .env.dist .env
    ``` 
1. Build and start the API containers
    - for testing and usage run
        ```bash
        $ make up
        ```
    - for development run
        ```bash
        $ make dev
        ```

The API should be accessible on the port number configured in `.env`,
e.g. [http://localhost:${API_PORT}](http://localhost:8080).


## API Specification

The OpenAPI Specification can be found in [doc/openapi.yml](doc/openapi.yml). 

To read the API specification and to test the API you can build and start the `doc` container.
```bash
$ make doc
```

This container will be accessible on the port number configured in your `.env` file as `DOC_PORT`,
e.g. [http://localhost:${DOC_PORT}](http://localhost:8081).