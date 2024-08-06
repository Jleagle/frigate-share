# frigate-share

    services:
        frigate-share:
            image: "ghcr.io/jleagle/frigate-share:main"
            container_name: "frigate-share"
            hostname: "frigate-share"
            restart: "unless-stopped"
            ports:
              - "5002:5002"
