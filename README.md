# SQL Processor CLI

This project provides a command-line tool for processing SQL files on a PostgreSQL database. The tool is packaged as a Docker container for easy deployment and usage.

## Usage of CLI 

1. Build the SQL Processor CLI:
    ```sh
    export SSL_MODE=disable
    cd cmd/sqlprocessorcli
    go build
    ```
2. Run `./sqlprocessorcli --help`
    - Export environment variables to configure the SQL Processor CLI:
       - **SSL_MODE** 
            - Options: `disable`, `allow`, `prefer`, `require`, `verify-ca`, `verify-full`. 
            - Default: `SSL_MODE=disable`.
       - **PG_HOST**
            - Set the PostgreSQL database connection. Options: `PG_HOST=host.docker.internal` (macOs & Windows). 
            - Default: `PG_HOST=localhost`.  
            

- Processes a single .sql file: `./sqlprocessorcli <file.sql>`
- Recursively processes .sql files at path: `./sqlprocessorcli </path/to/files/>` 


3. Dependencies:
    - [github.com/lib/pq](https://github.com/lib/pq)
    - [github.com/stretchr/testify](https://github.com/stretchr/testify)
	- [github.com/ScooterHelmet/procedures/pkg/sqlprocessor](https://github.com/ScooterHelmet/procedures)
    - [github.com/rubenv/sql-migrate](github.com/rubenv/sql-migrate)
    - [github.com/DATA-DOG/go-sqlmock](github.com/DATA-DOG/go-sqlmock)
	- [github.com/davecgh/go-spew](github.com/davecgh/go-spew)
	- [github.com/go-gorp/gorp/v3](github.com/go-gorp/gorp/v3)
	- [github.com/pmezard/go-difflib](github.com/pmezard/go-difflib)
    - [github.com/spf13/afero](github.com/spf13/afero)
    - [golang.org/x/text](golang.org/x/text)
	- [gopkg.in/yaml.v3](gopkg.in/yaml.v3)

4. Replace `yourusername`, `youruser`, `yourpassword`, and `yourdbname` with the appropriate values. The `SslMode` is set based on the environment variable `SSL_MODE`. If the system administrator does not provide the configuration, the default value is "disable", and a warning message is displayed.

## Usage of docker-compose.yml
The `docker-compose.yml` file defines a single service named sqlprocessor. It uses the sqlprocessor-cli Docker image and builds it from the current directory. The SSL_MODE environment variable is set to disable. The ./sqlfiles directory is mounted as a volume to the /sqlfiles directory inside the container.

1. Build the Docker image:
    ```bash
    docker-compose build
    ```
2. Create a `./sqlfiles` directory and place the .sql file you want to process.
3. Run the sqlprocessor service with Docker Compose:
    ```bash
    docker-compose up -d
    ```
4. The sqlprocessor service will process the `.sql` file in the `./sqlfiles` directory.

Configuration
You can configure the sqlprocessor service by setting environment variables in the docker-compose.yml file.

SSL_MODE: Set the SSL mode for the PostgreSQL connection. Options are disable, allow, prefer, require, verify-ca, and verify-full. The default is disable.

### Install and run a local PostgreSQL database

1. Linux (Ubuntu)
    -    Update the package lists and install the PostgreSQL server:
    ```sh
    sudo apt update
    sudo apt install postgresql postgresql-contrib
    ```
    -    Start the PostgreSQL service:
    ```sh
    sudo systemctl start postgresql
    ```
    -    Enable the PostgreSQL service to start on boot:
    ```sh
    sudo systemctl enable postgresql
    ```
    -    Change to the postgres user and access the PostgreSQL shell:
    ```sh
    sudo -u postgres psql
    ```
2. macOS
    -    Install the PostgreSQL server using Homebrew:
    ```sh
    brew install postgresql
    ```
    -    Start the PostgreSQL service:
    ```sh
    brew services start postgresql
    ```
    -   A homebrew install requires a default user:
    ```sh
    /usr/local/opt/postgresql\@14/bin/createuser -s postgres
    ```
    -    Access the PostgreSQL shell:
    ```sh
    psql postgres
    ```
    
3. Windows
    -    Download the Windows installer for PostgreSQL from the official website:
https://www.enterprisedb.com/downloads/postgres-postgresql-downloads
    -    Run the installer and follow the on-screen instructions to install PostgreSQL.
    -    After the installation is complete, the PostgreSQL service should start automatically. You can manage the service using the "pgAdmin" GUI tool that was installed with PostgreSQL.
    -    To access the PostgreSQL shell, open the "SQL Shell (psql)" application from the Start Menu.

### Setup a PostgreSQL database, user, & privileges via terminal
Regardless of the operating system, once you are in the PostgreSQL shell, you can create and manage databases, users, and execute SQL commands. For example, to create a new database and user, follow these steps:

1. Create a new database:
    ```sql
    CREATE DATABASE onet;
    ```

2. Create a new user with a password:
    ```sql
    CREATE USER onet WITH PASSWORD 'fathomrocks';
    ```

3. Grant privileges to the new user on the new database:
    ```sql
    GRANT ALL PRIVILEGES ON DATABASE onet TO onet;
    ```

4. Exit the PostgreSQL shell:
    ```psql
    \q
    ``` 

### Authenticate access to postgresql with pgAdmin

pgAdmin is a popular open-source administration and management tool for the PostgreSQL database. Here's how to connect to a PostgreSQL database using pgAdmin:

1. Install pgAdmin if you haven't already:

    -    **Linux (Ubuntu):** You can install pgAdmin using the following commands:
    ```sh
    sudo apt install pgadmin4
    ```  

    -    **macOS:** Install pgAdmin using Homebrew:
    ```sh
    brew install --cask pgadmin4
    ```  
    -    **Windows:** The pgAdmin installer is bundled with the PostgreSQL installer. If you've already installed PostgreSQL using the installer, you should have pgAdmin installed as well.

2. Launch pgAdmin:   
    -    **Linux:** Run pgadmin4 from the terminal or search for it in your applications menu.

    -    **macOS:** Search for "pgAdmin" in your applications, or run open `/Applications/pgAdmin\ 4.app` in the terminal.  

    -    **Windows:** Search for "pgAdmin" in the Start menu and open it.  

3. Add a new server connection:

    -    In the pgAdmin application window, go to the "Browser" panel on the left side, right-click "Servers" and select "Create" > "Server".  

    -    In the "Create - Server" window, switch to the "General" tab, and provide a name for your connection (e.g., "Local PostgreSQL").  

4. Switch to the "Connection" tab, and fill in the following fields:
    -    **Hostname/address:** Enter "localhost" if the PostgreSQL server is running on the same machine as pgAdmin, or provide the IP address or hostname of the remote server.
    -    **Port:** The default PostgreSQL port is 5432, but enter the appropriate port if it's different.
    -    **Maintenance database:** Enter "postgres" (the default maintenance database) or the name of the specific database you want to connect to.
    -    **Username:** Enter the username of the PostgreSQL user you want to use for this connection.
    -    **Password:** Enter the password for the PostgreSQL user.
    -    Check the "Save password" box if you want to save the password for future connections.
    -    Click the "Save" button to create the connection.

5. Explore the connected server:
   -    In the "Browser" panel, you should now see the new server connection listed under "Servers". Click the arrow next to it to expand the server's tree view.  
   -    Expand "Databases" to see a list of available databases. You can interact with the databases, tables, and other objects by right-clicking them and selecting the desired operations from the context menu.  

With pgAdmin, you can manage databases, tables, indexes, and other database objects, as well as run SQL queries, import/export data, and more.
