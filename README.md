# es_backend

## Image Storage Server (MinIO)

- To minimize costs, we utilized **MinIO** locally with Docker for our image
  storage server.
- Since our current internet setup doesn't allow for port forwarding, and with a
  small development team, we resolved the port forwarding issue by using
  **Tailscale**.

## Character Classification

- Characters are categorized based on their weapon types and their unique
  combinations.

  - The classification is based on a specific position list, referenced from
    Youtuber 모좀's tier list:

    - Tank (탱커)

    - Tank Bruiser (탱 브루저)

    - Damage Bruiser (딜 브루저)

    - Assassin (암살자)

    - Auto Attack Dealer (평원딜)

    - Skill Amplification Mage (스증 마법사)

    - Support (서포터)
