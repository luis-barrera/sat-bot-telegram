# This snippets allows to execute this app as a daemon
# But, first clone the code and build a image of the root project with the
#  command: docker build -t sat-bot-telegram .
{
  # NOTE: make the corresponding modifications to your orignal nixos config
  systemd.services.sat-bot = {
    description = "SAT Bot Docker Container";
    after = [ "network.target" "docker.service" ];
    wants = [ "docker.service" ];  # Requires Docker to be running
    wantedBy = [ "multi-user.target" ];  # Starts at boot

    serviceConfig = {
      ExecStart = "${pkgs.docker}/bin/docker start -a sat-bot";
      ExecStop = "${pkgs.docker}/bin/docker stop -t 10 sat-bot";
      Restart = "on-failure";  # Restarts on non-clean exit
      RestartSec = "10s";      # Wait 10 seconds before restarting
    };

    # This ensures the container exists before starting
    preStart = ''
      if ! ${pkgs.docker}/bin/docker container inspect sat-bot >/dev/null 2>&1; then
        ${pkgs.docker}/bin/docker run -d --name sat-bot sat-bot-telegram
      fi
    '';
  };
}
