# NAME

cskr - Curvy Social Key Recovery

# SYNOPSIS

cskr

**Usage**:

```
cskr [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# COMMANDS

## split, s

Generate the shares for the private (k,v) pair

**--nOfShares, -n**="": Number of shares (default: 0)

**--output, -o**="": Output directory (default: .)

**--threshold, -t**="": Number of shares required to reconstruct the (k,v) pair (default: 0)

## recover, r

Recover the private (k,v) pair from the given shares

**--output, -o**="": Output path for (k,v) JSON file

**--pattern, -p**="": Pattern for matching share files (default: share*.json)

**--threshold, -t**="": Number of shares needed for recovery (default: 0)
