{
  config,
  pkgs,
  lib,
  ...
}: {
  options = with lib;
  with types; {
    services.dynamic-ip-watcher = {
      enable = mkOption {
        type = types.bool;
        default = false;
        description = ''
          Enable Dynamic IP Watcher.
        '';
      };
      interval = mkOption {
        type = str;
        default = "1m";
        description = "Interval at which to run the Dynamic IP Watcher service (e.g., '1h', '30m').";
      };
      dnsRecord = mkOption {
        description = "Options for a managed DNS Record.";
        default = {};
        type = submodule {
          options = {
            type = mkOption {
              type = enum ["cloudflare" "none"];
              default = "none";
              description = "Type of DNS provider.";
            };
            apiKey = mkOption {
              type = str;
              default = "";
              description = "API key for the DNS provider.";
            };
            zoneName = mkOption {
              type = str;
              default = "";
              description = "Zone name for the DNS provider.";
            };
            recordName = mkOption {
              type = str;
              default = "";
              description = "Record name for the DNS provider.";
            };
          };
        };
      };
      storage = mkOption {
        description = "Options for storage.";
        default = {};
        type = submodule {
          options = {
            directory = mkOption {
              type = str;
              default = "/var/lib/dynamic-ip-watcher";
              description = "The directory the application will store data.";
            };
          };
        };
      };
      notifiers = mkOption {
        description = "Endpoints to notify on change.";
        default = [];
        type = listOf (submodule {
          options = {
            type = mkOption {
              type = enum ["discord"];
              description = ''
                The type of notifier.
              '';
            };
            webhookUrl = mkOption {
              type = str;
              default = "";
              description = ''
                The webhook url to send the message.
              '';
            };
            username = mkOption {
              type = str;
              default = "";
              description = ''
                The username to use in the message. Used with 'discord' type.
              '';
            };
            avatarUrl = mkOption {
              type = str;
              default = "";
              description = ''
                The avatar url to use in the message. Used with 'discord' type.
              '';
            };
          };
        });
      };
    };
  };
}
