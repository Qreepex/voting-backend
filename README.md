# voting-backend
🧮 Backend for Voting Web App


## Routes

### POST /v1/vote
Body:
```json
{
    "candidate": "candidate_id",
    "campaign": "campaign_id", // to track where the vote is coming from
}
```

### GET /v1/votes/{candidate_id}
Response:
```json
[
    voteTimestamp1,
    voteTimestamp2,
    ...
]
```

### GET /v1/candidates
Response:
```json
[
    {
        "name": "candidate_name",
        "id": "candidate_id",
        ...other_candidate_data
    }
]

```

### GET /v1/candidate/{candidate_id}
Response:
```json
{
    "name": "candidate_name",
    "id": "candidate_id",
    ...other_candidate_data
}
```

## 24h Vote Cooldown Strategy

### Cookie based

- Check if Request has a cookie with name COOKIE_NAME
- If yes, check if value is in Redis
- If yes, return 403
- Check if IP is in Redis
- If yes, return 403
- If no, add cookie and ip to Redis and set cookie
- If cookie is not present, add value to Redis and set cookie without TTL

### IP based (Fallback)

- Check if Request IP is in Redis
- If yes, return 403
- If no, add IP to Redis with TTL


- Kein Cookie, nach IP gucken, wenn IP, den Cookie nehmen und http setzrn