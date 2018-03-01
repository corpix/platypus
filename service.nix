{ config, lib, pkgs, ... }:

with lib;

let
  name = "platypus";
  cfg = config.services."${name}";
  pkg = (pkgs.callPackage ./default.nix { }).bin;
in {
  options = with types; {
    services."${name}" = {
      enable = mkEnableOption "Platypus HTTP+WEBSOCKET data server";
      user = mkOption {
        default = name;
        type = string;
        description = ''
          User name to run service from.
        '';
      };
      group = mkOption {
        default = name;
        type = string;
        description =''
          Group name to run service from.
        '';
      };
      configuration = mkOption {
        default = {};
        description = ''
          Application configuration.
        '';
      };
    };
  };

  config = mkIf cfg.enable {
    users.extraUsers."${name}" = {
      name = name;
      group = cfg.group;
      uid = config.cryptounicorns.ids.uids."${name}";
    };

    users.extraGroups."${name}" = {
      name = name;
      gid = config.cryptounicorns.ids.gids."${name}";
    };

    systemd.services."${name}" = {
      enable = true;

      wants    = [ "nsqd.service" "nginx.service"  ];
      wantedBy = [ "multi-user.target" ];
      after    = [ "network.target" ];

      serviceConfig = {
        Type = "simple";
        User = name;
        Group = name;
        ExecStart = "${pkg}/bin/${name} -c ${pkgs.writeText "config.json" (builtins.toJSON cfg.configuration)}";
        Restart = "on-failure";
        RestartSec = 1;
      };
    };
  };
}
