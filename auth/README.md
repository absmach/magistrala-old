# Auth - Authentication and Authorization service

Auth service provides authentication features as an API for managing authentication keys as well as administering groups of entities - `things` and `users`.

## Authentication

User service is using Auth service gRPC API to obtain login token or password reset token. Authentication key consists of the following fields:

- ID - key ID
- Type - one of the three types described below
- IssuerID - an ID of the Magistrala User who issued the key
- Subject - user email
- IssuedAt - the timestamp when the key is issued
- ExpiresAt - the timestamp after which the key is invalid

There are _three types of authentication keys_:

- User key - keys issued to the user upon login request
- API key - keys issued upon the user request
- Recovery key - password recovery key

Authentication keys are represented and distributed by the corresponding [JWT](jwt.io).

User keys are issued when user logs in. Each user request (other than `registration` and `login`) contains user key that is used to authenticate the user.

API keys are similar to the User keys. The main difference is that API keys have configurable expiration time. If no time is set, the key will never expire. For that reason, API keys are _the only key type that can be revoked_. This also means that, despite being used as a JWT, it requires a query to the database to validate the API key. The user with API key can perform all the same actions as the user with login key (can act on behalf of the user for Thing, Channel, or user profile management), _except issuing new API keys_.

Recovery key is the password recovery key. It's short-lived token used for password recovery process.

For in-depth explanation of the aforementioned scenarios, as well as thorough
understanding of Magistrala, please check out the [official documentation][doc].

The following actions are supported:

- create (all key types)
- verify (all key types)
- obtain (API keys only)
- revoke (API keys only)

## Groups

User and Things service are using Auth gRPC API to get the list of ids that are part of a group. Groups can be organized as tree structure.
Group consists of the following fields:

- ID - ULID id uniquely representing group
- Name - name of the group, name of the group is unique at the same level of tree hierarchy for a given tree.
- ParentID - id of the parent group
- OwnerID - id of the user that created a group
- Description - free form text, up to 1024 characters
- Metadata - Arbitrary, object-encoded group's data
- Path - tree path consisting of group ids
- CreatedAt - timestamp at which the group is created
- UpdatedAt - timestamp at which the group is updated

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                     | Description                                                             | Default                          |
| ---------------------------- | ----------------------------------------------------------------------- | -------------------------------- |
| MG_AUTH_LOG_LEVEL            | Service level (debug, info, warn, error)                                | error                            |
| MG_AUTH_DB_HOST              | Database host address                                                   | localhost                        |
| MG_AUTH_DB_PORT              | Database host port                                                      | 5432                             |
| MG_AUTH_DB_USER              | Database user                                                           | magistrala                       |
| MG_AUTH_DB_PASSWORD          | Database password                                                       | magistrala                       |
| MG_AUTH_DB                   | Name of the database used by the service                                | auth                             |
| MG_AUTH_DB_SSL_MODE          | Database connection SSL mode (disable, require, verify-ca, verify-full) | disable                          |
| MG_AUTH_DB_SSL_CERT          | Path to the PEM encoded certificate file                                |                                  |
| MG_AUTH_DB_SSL_KEY           | Path to the PEM encoded key file                                        |                                  |
| MG_AUTH_DB_SSL_ROOT_CERT     | Path to the PEM encoded root certificate file                           |                                  |
| MG_AUTH_HTTP_PORT            | Auth service HTTP port                                                  | 8180                             |
| MG_AUTH_GRPC_PORT            | Auth service gRPC port                                                  | 8181                             |
| MG_AUTH_SERVER_CERT          | Path to server certificate in pem format                                |                                  |
| MG_AUTH_SERVER_KEY           | Path to server key in pem format                                        |                                  |
| MG_AUTH_SECRET               | String used for signing tokens                                          | auth                             |
| MG_AUTH_LOGIN_TOKEN_DURATION | The login token expiration period                                       | 10h                              |
| MG_JAEGER_URL                | Jaeger server URL                                                       | <http://jaeger:14268/api/traces> |
| MG_KETO_READ_REMOTE_HOST     | Keto Read Host                                                          | magistrala-keto                  |
| MG_KETO_WRITE_REMOTE_HOST    | Keto Write Host                                                         | magistrala-keto                  |
| MG_KETO_READ_REMOTE_PORT     | Keto Read Port                                                          | 4466                             |
| MG_KETO_WRITE_REMOTE_PORT    | Keto Write Port                                                         | 4467                             |

## Deployment

The service itself is distributed as Docker container. Check the [`auth`](https://github.com/absmach/magistrala/blob/master/docker/docker-compose.yml#L71-L94) service section in
docker-compose to see how service is deployed.

To start the service outside of the container, execute the following shell script:

```bash
# download the latest version of the service
go get github.com/absmach/magistrala

cd $GOPATH/src/github.com/absmach/magistrala

# compile the service
make auth

# copy binary to bin
make install

# set the environment variables and run the service
MG_AUTH_LOG_LEVEL=[Service log level] MG_AUTH_DB_HOST=[Database host address] MG_AUTH_DB_PORT=[Database host port] MG_AUTH_DB_USER=[Database user] MG_AUTH_DB_PASS=[Database password] MG_AUTH_DB=[Name of the database used by the service] MG_AUTH_DB_SSL_MODE=[SSL mode to connect to the database with] MG_AUTH_DB_SSL_CERT=[Path to the PEM encoded certificate file] MG_AUTH_DB_SSL_KEY=[Path to the PEM encoded key file] MG_AUTH_DB_SSL_ROOT_CERT=[Path to the PEM encoded root certificate file] MG_AUTH_HTTP_PORT=[Service HTTP port] MG_AUTH_GRPC_PORT=[Service gRPC port] MG_AUTH_SECRET=[String used for signing tokens] MG_AUTH_SERVER_CERT=[Path to server certificate] MG_AUTH_SERVER_KEY=[Path to server key] MG_JAEGER_URL=[Jaeger server URL] MG_AUTH_LOGIN_TOKEN_DURATION=[The login token expiration period] $GOBIN/magistrala-auth
```

If `MG_EMAIL_TEMPLATE` doesn't point to any file service will function but password reset functionality will not work.

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](https://api.mainflux.io/?urls.primaryName=auth-openapi.yml).

[doc]: https://docs.mainflux.io
