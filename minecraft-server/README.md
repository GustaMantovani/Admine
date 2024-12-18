# Minecraft Server with ZeroTier

The Minecraft Server with ZeroTier is a Docker-based solution for running a Minecraft server with VPN connectivity provided by ZeroTier. This setup ensures secure and consistent access to the Minecraft server, leveraging Docker Compose for easy orchestration and management.

## Environment Variables

The following environment variables are used to configure the Docker build and runtime environment:

```bash
export JAVA_VERSION=17
export NETWORK_ID=a123b456c789
```

## Configuration Files

Configuration files for the Minecraft server should be placed in the `config` directory. These files will be mounted as volumes in the Docker container to allow for easy modification and persistence.

### Configuration Files

- `config/eula.txt`: Contains the EULA agreement for the Minecraft server.
- `config/server.properties`: Contains the server properties configuration.
- `config/user_jvm_args.txt`: Contains JVM arguments for the Minecraft server.
