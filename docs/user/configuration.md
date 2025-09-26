# Configuration

The k2ray application can be configured using environment variables or a `.env` file. The application will look for a file named `system.env` in the `configs/` directory by default.

## Configuration Method

You can either create a `configs/system.env` file in the root of the project or set the environment variables directly in your shell or deployment environment. Environment variables will always take precedence over values defined in the `.env` file.

### Example `configs/system.env` file:

```env
# The connection string for the SQLite database.
# Defaults to a local file 'k2ray.db' if not set.
DATABASE_URL=./k2ray.db

# A secret key for signing JWT (JSON Web Tokens).
# It is strongly recommended to change this to a long, random string in production.
# Defaults to 'default-secret-please-change'.
JWT_SECRET=your-super-secret-and-long-random-string
```

## Configuration Variables

Here are the environment variables used by the application:

| Variable      | Description                                                                                                | Default Value                    |
|---------------|------------------------------------------------------------------------------------------------------------|----------------------------------|
| `DATABASE_URL`| The connection string for the database. For SQLite, this is typically the path to the database file.         | `./k2ray.db`                     |
| `JWT_SECRET`  | The secret key used to sign and verify JSON Web Tokens (JWTs) for user authentication.                       | `default-secret-please-change`   |

**Important:** For a production environment, always set a strong, unique `JWT_SECRET` to ensure the security of your application. Do not use the default value.