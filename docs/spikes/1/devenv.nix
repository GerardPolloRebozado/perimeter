{ pkgs, ... }:

{
  # https://devenv.sh/basics/
  env.GREET = "Loco Development Environment";

  # Enable PostgreSQL service
  services.postgres = {
    enable = true;
    package = pkgs.postgresql_15; # You can match your preferred version

    # Automatically create the database your Loco app expects on startup
    initialDatabases = [
      { name = "surface_development"; }
    ];
    initialScript = ''
      CREATE USER loco WITH SUPERUSER PASSWORD 'loco';
    '';
    # Optional: ensure it listens on standard TCP port so connection strings work seamlessly
    listen_addresses = "127.0.0.1";
  };

  # See full list of options at https://devenv.sh/reference/options/
}
