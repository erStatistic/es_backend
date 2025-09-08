# es_backend

## Image Storage Server (MinIO)

- To minimize costs, we utilized **MinIO** locally with Docker for our image
  storage server.
- Since our current internet setup doesn't allow for port forwarding, and with a
  small development team, we resolved the port forwarding issue by using
  **Tailscale**.

## Character Cluster

- To cluster characters with weapon, we choose the way called Kmeans.
- Actually, We can choose others but, the value of Kmeans is smiliar to our real
  value.

## Character Scoring

- To calculate a score for each character, the following steps are applied:

  - Collect Win Rates by Combination
    - For each character, gather win rates for all weapon combinations they can
      use.
  - Apply Weight Based on Sample Size
    - Weight each win rate by the number of matches played with that
      combination.
    - This ensures that combinations with more data have a stronger influence on
      the score.
  - Average Weighted Scores
    - Calculate the weighted average win rate across all combinations for the
      character.
  - Normalize Scores
    - Normalize the results so that scores can be compared fairly across all
      characters.
  - Optional: Role-Based Adjustment

- Apply additional weighting or adjustments based on the character’s intended
  role
  - (e.g., Support, DPS, Tank).
