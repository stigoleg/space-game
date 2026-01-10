# GitHub Gist Online Leaderboard Setup Guide

## Overview

Stellar Siege supports an **optional online leaderboard** using GitHub Gist as the backend storage. This is completely free and requires no server setup!

## How It Works

1. **GitHub Gist Storage**: Your leaderboard data is stored in a public GitHub Gist (JSON format)
2. **API Access**: The game reads and writes scores via GitHub's Gist API
3. **Local Caching**: Scores are cached for 30 seconds to reduce API calls
4. **Async Submission**: Score submissions happen in the background without blocking gameplay

## Setup Instructions (5 minutes)

### Step 1: Create a GitHub Account (if you don't have one)
- Go to https://github.com/signup
- Create your free account

### Step 2: Create a GitHub Personal Access Token

1. Go to https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Name it: `stellar-siege-leaderboard`
4. Check only this permission: `gist` (read/write gists)
5. Set expiration as desired (or "No expiration")
6. Click "Generate token"
7. **Copy the token and save it somewhere safe** (you won't see it again!)

Example token format: `ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`

### Step 3: Create a Gist for Your Leaderboard

1. Go to https://gist.github.com
2. Filename: `leaderboard.json`
3. Content (paste this):
```json
[]
```
4. Select "Create public gist"
5. Copy the Gist ID from the URL

Example URL: `https://gist.github.com/yourname/a1b2c3d4e5f6`  
The Gist ID is: `a1b2c3d4e5f6`

### Step 4: Configure Stellar Siege

Create a `.env` file in the game directory (next to the executable):

```env
GIST_ID=your_gist_id_here
GH_GIST_TOKEN=your_github_token_here
GIST_ENABLED=true
```

Example:
```env
GIST_ID=a1b2c3d4e5f6
GH_GIST_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
GIST_ENABLED=true
```

**Important**: 
- Use `GH_GIST_TOKEN`, not `GITHUB_TOKEN`
- Never commit the `.env` file to a public repository
- Keep your token private - it grants write access to your gist

### Step 5: Run the Game

```bash
./stellar-siege
```

That's it! When you complete a game, you'll be prompted to submit your score to the online leaderboard.

---

## Features

### Score Submission
- **Automatic prompt** after each game
- **Player name** is required
- **Difficulty level** is recorded (Easy/Normal/Hard)
- **Wave reached** is recorded
- **Timestamp** is automatically added

### Online Leaderboard Display
- **Top 10 global scores** displayed after game over
- **Your score highlighted** if you're on the leaderboard
- Shows: Rank, Player Name, Score, Difficulty, Wave Reached

### Local Leaderboard
- Separate local leaderboard for your device
- 30-second cache for fast display
- Automatically syncs with online leaderboard

### Data Storage
All scores are stored as JSON in the Gist:
```json
[
  {
    "player_name": "Alice",
    "score": 15000,
    "difficulty": "Hard",
    "date": "2026-01-04T13:00:00Z",
    "wave": 12
  }
]
```

---

## Troubleshooting

### Score Not Submitting?

**Problem**: "Online leaderboard not configured" message
- **Solution**: Check that `.env` file exists with correct values
- **Verify**: 
  - Is `GIST_ENABLED` set to `true`?
  - Is `GIST_ID` not empty?
  - Is `GH_GIST_TOKEN` valid?

**Problem**: API returns 401 Unauthorized
- **Solution**: Your GitHub token is invalid or expired
- **Fix**: Generate a new token at https://github.com/settings/tokens

**Problem**: API returns 404 Not Found
- **Solution**: Your Gist ID is wrong
- **Fix**: Copy the correct ID from your Gist URL

**Problem**: Gist update fails silently
- **Solution**: Check your token has `gist` permission
- **Fix**: Create a new token with proper scopes

### Leaderboard Not Showing?

1. Make sure at least one score has been submitted
2. Wait 30 seconds (cache expiry) and check again
3. Verify the Gist is public: https://gist.github.com/yourname/GIST_ID

### Network Issues?

- Game works offline! Scores won't submit, but the game continues
- Errors are silent - if submission fails, you can retry next game

---

## Configuration Options

### Disable Online Leaderboard (Keep Local Only)

Edit your `.env` file:
```env
GIST_ENABLED=false
```

Or simply delete the `.env` file.

### Share Your Leaderboard URL

Once configured, your leaderboard is public at:
```
https://gist.githubusercontent.com/yourname/GIST_ID/raw/leaderboard.json
```

You can share this URL with friends!

### Leaderboard Limits

- **Maximum scores stored**: 100 (keeps only top 100)
- **Cache duration**: 30 seconds
- **API rate limit**: 5,000 requests/hour per token
- **Update latency**: Usually 500ms-2 seconds

---

## API Details (For Developers)

### GitHub Gist API Endpoints

**Read (public, no auth needed)**:
```
GET https://gist.githubusercontent.com/raw/{gist_id}/leaderboard.json
```

**Write (requires auth token)**:
```
PATCH https://api.github.com/gists/{gist_id}
Authorization: token YOUR_TOKEN
Content-Type: application/json
```

### Implementation Files

- **Backend**: `game/systems/gist_leaderboard.go`
- **Config**: `game/systems/gist_config.go`
- **Integration**: `game/game.go` (updateGameOver, submitScoreOnline, drawOnlineLeaderboard)

### Environment Variables

The game reads configuration from environment variables in this order:

1. **Environment variables** (highest priority)
   - `GIST_ID`
   - `GH_GIST_TOKEN`
   - `GIST_ENABLED`

2. **.env file** (loaded via godotenv)
   - Automatically loaded at startup
   - Place in same directory as executable

3. **JSON config file** (fallback, not recommended for secrets)
   - `config/gist_config.json`
   - **Warning**: Storing secrets in JSON is not secure!

**Recommendation**: Always use `.env` file for configuration, never commit secrets to your repository.

---

## Privacy & Security

### Data Privacy
- Leaderboard is **public by default** (visible to anyone with the Gist URL)
- Player names are visible
- Scores and timestamps are visible

### Data Safety
- Stored on GitHub's secure servers
- Automatic version history (can revert changes)
- GitHub has 99.99% uptime

### Token Security
- **Never share your token** - it grants write access to your gist
- **Never commit `.env` to a repository** - add it to `.gitignore`
- Token is read from `.env` file at startup
- Revoke and regenerate token if compromised

### Cheating Prevention
- Current implementation uses client-submitted scores
- To prevent cheating, you could:
  - Add server-side validation
  - Hash scores with timestamp
  - Implement anti-tampering checksums
  - Monitor for suspicious score patterns

---

## FAQ

**Q: Is it really free?**  
A: Yes! GitHub Gist is free, and the API has generous rate limits for personal use.

**Q: Can other people modify my leaderboard?**  
A: No! Only you can modify it with your token. Others can read it if they know the Gist ID.

**Q: What if I lose my token?**  
A: Generate a new one at https://github.com/settings/tokens. Old tokens become invalid.

**Q: Can I reset my leaderboard?**  
A: Yes! Edit the Gist directly at https://gist.github.com/yourname/GIST_ID and replace content with `[]`

**Q: How do I share my leaderboard?**  
A: Share your Gist URL: `https://gist.github.com/yourname/GIST_ID`

**Q: Can I see other players' leaderboards?**  
A: Yes! If they share their Gist URL, you can see their scores.

**Q: Why can't I use GITHUB_TOKEN as the secret name?**  
A: GitHub reserves secret names starting with `GITHUB_` for internal use. Use `GH_GIST_TOKEN` instead.

**Q: Should I store secrets in config/gist_config.json?**  
A: No! Always use `.env` file for secrets. JSON config is only for fallback and should not contain tokens.

---

## Support

If you have issues:

1. **Check the .env file** - most problems are configuration
2. **Verify your token** - test it at https://api.github.com/user (add header: `Authorization: token YOUR_TOKEN`)
3. **Check your Gist** - view it at https://gist.github.com/yourname/GIST_ID
4. **View the logs** - check game console for error messages

See [LEADERBOARD_CONFIG_FIX.md](LEADERBOARD_CONFIG_FIX.md) for detailed troubleshooting.

---

**Enjoy competing on the leaderboard!** 🏆
