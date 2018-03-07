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
      domain = mkOption {
        type = str;
        description = ''
          Domain which should be used for this service.
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

      wantedBy = [ "multi-user.target" ];
      after    = [ "nsqd.service" "network.target" ];

      serviceConfig = {
        Type = "simple";
        User = name;
        Group = name;
        ExecStart = "${pkg}/bin/${name} -c ${pkgs.writeText "config.json" (builtins.toJSON cfg.configuration)}";
        Restart = "on-failure";
        RestartSec = 1;
      };
    };

    services.nginx = {
      upstreams = {
        platypus = {
          servers = {
            "${cfg.configuration.HTTP.Addr}" = { backup = false; };
          };
        };
      };

      virtualHosts."${cfg.domain}".extraConfig = let
        proxy = path: ''
          location ${path.path} {
            proxy_pass         http://platypus;
            proxy_set_header   X-Forwarded-For $remote_addr;
            proxy_http_version 1.1;
            ${optionalString (path.type == "stream" || path.type == "streams") ''
              proxy_set_header   Upgrade         $http_upgrade;
              proxy_set_header   Connection      "upgrade";
            ''}
          }
        '';
        paths = configuration: map
          (handler: { path = handler.Path; type = handler.Type; })
          configuration.Handlers;
      in mkAfter (concatStringsSep "\n" (map proxy (paths cfg.configuration)));
    };
  };
}
