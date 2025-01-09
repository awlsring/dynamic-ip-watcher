overlay: {
  config,
  pkgs,
  lib,
  ...
}: let
  cfg = config.services.dynamic-ip-watcher;
  format = pkgs.formats.json {};
  filterNulls = lib.filterAttrsRecursive (v: v != null);
  configFile = format.generate "dynamic-ip-watcher.json" cfg;
in {
  imports = [./options.nix];
  
  config = lib.mkIf cfg.enable {
    nixpkgs.overlays = [overlay];
    environment.systemPackages = [pkgs.dynamic-ip-watcher];
    systemd.services.dynamic-ip-watcher = {
      description = "Dynamic IP Watcher Service";
      wantedBy = ["multi-user.target"];
      after = ["network-online.target"];
      wants = ["network-online.target"];

      serviceConfig = {
        Type = "oneshot";
        ExecStart = "${pkgs.dynamic-ip-watcher}/bin/dynamic-ip-watcher --config-path ${configFile}";
        User = "dynamic-ip-watcher";
        Group = "dynamic-ip-watcher";
        StateDirectory = "dynamic-ip-watcher";
        RemainAfterExit = false;
        Restart = "on-failure";
      };
    };

    systemd.timers.dynamic-ip-watcher = {
      description = "Timer for Dynamic IP Watcher Service";
      wantedBy = ["timers.target"];

      timerConfig = {
        OnUnitActiveSec = cfg.interval;
        OnBootSec = "1m";
        Unit = "dynamic-ip-watcher.service";
      };
    };

    users.groups.dynamic-ip-watcher = {};
    users.users.dynamic-ip-watcher = {
      description = "Dynamic IP Watcher service user";
      group = "dynamic-ip-watcher";
      isSystemUser = true;
    };
  };
}
